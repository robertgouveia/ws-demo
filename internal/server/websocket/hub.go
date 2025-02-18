package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Hub struct {
	clients     map[uint][]*Client
	lock        sync.RWMutex
	upgrader    *websocket.Upgrader
	errLogger   *log.Logger
	infoLogger  *log.Logger
	pong        time.Duration
	ping        time.Duration
	gameManager *GameManager
}

func NewHub(errLogger *log.Logger, infoLogger *log.Logger, gm *GameManager) *Hub {
	return &Hub{
		clients: make(map[uint][]*Client),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		errLogger:   errLogger,
		infoLogger:  infoLogger,
		pong:        time.Second * 5,
		ping:        time.Second * 2,
		gameManager: gm,
	}
}

func (h *Hub) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.errLogger.Printf("Unable to upgrade connection: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	h.infoLogger.Printf("New connection from %s", conn.RemoteAddr())

	client := NewClient(conn, uint(len(h.clients)+1), h.RemoveClient, h.infoLogger, h.errLogger)

	h.gameManager.Join(client)
}

func (h *Hub) registerClient(userID uint, client *Client) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[userID] = append(h.clients[userID], client)
}

func (h *Hub) RemoveClient(client *Client) {
	h.infoLogger.Printf("Removing client from hub: %s", client.conn.RemoteAddr())
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.clients, client.userId)
}
