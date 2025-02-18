package websocket

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"sync"
)

type game struct {
	ID      string
	Player1 *Client
	Player2 *Client
	Round   int
	Status  string
}

type GameManager struct {
	games    []*game
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

			err, code := m.NewGame(client1, client2)
			if err != nil {
				m.errorLog.Println(err)
			} else {
				log.Printf("Game Created under ID: %s", code)
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

func (m *GameManager) NewGame(p1, p2 *Client) (error, string) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return err, ""
	}
	code := hex.EncodeToString(bytes)

	game := &game{
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

	return nil, code
}
