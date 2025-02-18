package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Client struct {
	userId   uint
	conn     *websocket.Conn
	remove   func(*Client)
	send     chan []byte
	infoLog  *log.Logger
	errorLog *log.Logger
	status   string
	lock     sync.RWMutex
}

func NewClient(conn *websocket.Conn, userId uint, remove func(client *Client), infoLog *log.Logger, errorLog *log.Logger) *Client {
	return &Client{
		conn:     conn,
		userId:   userId,
		remove:   remove,
		send:     make(chan []byte, 256),
		infoLog:  infoLog,
		errorLog: errorLog,
		status:   "matching",
	}
}

func (c *Client) ChangeStatus(status string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.status = status
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Read() {
	defer func() {
		c.infoLog.Println("Client disconnected: %s", c.conn.RemoteAddr())
		c.remove(c)
		_ = c.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.infoLog.Println(string(message))
	}
}
