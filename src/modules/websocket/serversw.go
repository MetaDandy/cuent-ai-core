package ws

import "github.com/gofiber/websocket/v2"

// ServeWs devuelve un handler para Fiber/WebSocket que gestiona salas.
func ServeWs(h *Hub) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		roomID := c.Query("room")
		userID := c.Query("user")
		client := NewClient(roomID, userID, c)
		go client.WritePump()
		if err := h.Register(roomID, userID, client); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("sala llena"))
			c.Close()
			return
		}
		client.ReadPump(h)
	}
}
