package websocket

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Game struct {
	ID      string
	Player1 *Client
	Player2 *Client
	Round   int
	Status  string
	lock    sync.RWMutex
}

type GameManager struct {
	games    []*Game
	lock     sync.RWMutex
	infoLog  *log.Logger
	errorLog *log.Logger
	joiners  []*Client
}

func NewGameManager(ilog, elog *log.Logger) *GameManager {
	return &GameManager{
		joiners:  make([]*Client, 0),
		errorLog: elog,
		infoLog:  ilog,
	}
}

func (m *GameManager) ListenForJoin() {
	for {
		m.lock.Lock()
		if len(m.joiners) >= 2 {
			client1 := m.joiners[0]
			client2 := m.joiners[1]

			m.joiners = m.joiners[2:]
			m.lock.Unlock()

			err, game, code := m.NewGame(client1, client2)
			if err != nil {
				m.errorLog.Println(err)
			} else {
				log.Printf("Game Created under ID: %s", code)
				go game.StartGame()
			}
		} else {
			m.lock.Unlock()
		}
	}
}

func (m *GameManager) Join(player *Client) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.joiners = append(m.joiners, player)
	log.Printf("Player joined: %v", player)
}

func (m *GameManager) NewGame(p1, p2 *Client) (error, *Game, string) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return err, nil, ""
	}
	code := hex.EncodeToString(bytes)

	game := &Game{
		ID:      code,
		Player1: p1,
		Player2: p2,
		Round:   1,
		Status:  "pending",
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.games = append(m.games, game)

	m.infoLog.Printf("New game: %s", game.ID)

	p1.ChangeStatus("matched")
	p2.ChangeStatus("matched")

	return nil, game, code
}

func (g *Game) broadcast(message string, conn ...*websocket.Conn) {
	for _, conn := range conn {
		g.send(conn, message)
	}
}

func (g *Game) send(conn *websocket.Conn, message string) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (g *Game) StartGame() {
	g.lock.Lock()
	g.Status = "active"
	g.lock.Unlock()
	g.broadcast("Match Found!", g.Player1.conn, g.Player2.conn)

	// Game Loop
	for g.Status == "active" {
		g.broadcast("Player 1's turn...", g.Player1.conn, g.Player2.conn)
		g.Turn(g.Player1.conn, "Player 1")

		g.broadcast("Player 2's turn...", g.Player1.conn, g.Player2.conn)
		g.Turn(g.Player2.conn, "Player 2")

		if g.Status == "cancelled" || g.Status == "finished" {
			break
		}
	}

	g.broadcast("Game Finished!", g.Player1.conn, g.Player2.conn)
}

func (g *Game) Turn(conn *websocket.Conn, name string) {
	for {
		g.send(conn, "Please pick a tile")
		_, message, err := conn.ReadMessage()
		if err != nil {
			g.lock.Lock()
			g.Status = "cancelled"
			g.lock.Unlock()
			return
		}

		hit := string(message)

		g.broadcast(fmt.Sprintf("%s hit %s", name, hit), g.Player1.conn, g.Player2.conn)
		break
	}
}
