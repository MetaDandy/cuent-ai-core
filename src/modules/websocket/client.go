package ws

import "github.com/gofiber/websocket/v2"

// Client representa una conexión en una sala.
type Client struct {
	roomID string
	userID string
	conn   *websocket.Conn
	send   chan []byte
}

// NewClient crea un cliente para la sala y usuario especificados.
func NewClient(roomID, userID string, conn *websocket.Conn) *Client {
	return &Client{roomID: roomID, userID: userID, conn: conn, send: make(chan []byte, 256)}
}

// ReadPump lee mensajes del cliente y los retransmite a la sala.
func (c *Client) ReadPump(h *Hub) {
	defer h.Unregister(c.roomID, c.userID)
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		if string(msg) == `{"type":"disconnect"}` {
			return // sale de ReadPump y dispara defer Unregister
		}

		h.Broadcast(c.roomID, msg)
	}
}

// WritePump envía mensajes desde el canal send al cliente.
func (c *Client) WritePump() {
	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}
