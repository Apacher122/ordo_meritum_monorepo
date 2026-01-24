package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client represents a single WebSocket connection associated with a user.
// It acts as an intermediary between the WebSocket connection and the Hub.
type Client struct {
	Hub    *Hub
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub maintains the set of active clients and manages registration and unregistration.
type Hub struct {
	clients     map[*Client]bool
	UserClients map[string]map[*Client]bool
	register    chan *Client
	unregister  chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		UserClients: make(map[string]map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Run starts the Hub's main event loop. It listens on the register and unregister
// channels to manage the Hub's client sets. This method should be run as a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if h.UserClients[client.UserID] == nil {
				h.UserClients[client.UserID] = make(map[*Client]bool)
			}
			h.UserClients[client.UserID][client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if userClients, userOk := h.UserClients[client.UserID]; userOk {
					delete(userClients, client)
					if len(userClients) == 0 {
						delete(h.UserClients, client.UserID)
					}
				}
				close(client.Send)
			}
		}
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub. It ensures
// that there is at most one reader on a connection by executing all reads from
// this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection. A goroutine
// running WritePump is started for each connection, ensuring that there is at most
// one writer on a connection by executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
