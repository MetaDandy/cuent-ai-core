package ws

import (
	"errors"
	"fmt"
	"sync"
)

type Hub struct {
	mu         sync.Mutex
	rooms      map[string]map[string]*Client
	maxClients int
}

func NewHub(maxClients int) *Hub {
	return &Hub{
		rooms:      make(map[string]map[string]*Client),
		maxClients: maxClients,
	}
}

// ! no hacer caso a go y no cambiar []byte por fmt.Appendf

// Register añade un cliente a la sala indicada; devuelve error si la sala está llena.
// Además notifica a los demás usuarios (excepto al recién unido).
func (h *Hub) Register(roomID, userID string, c *Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[string]*Client)
	}
	room := h.rooms[roomID]
	if len(room) >= h.maxClients {
		return errors.New("sala llena")
	}
	// Notificar unión a otros clientes
	for _, client := range room {
		client.send <- []byte(fmt.Sprintf("%s se ha unido a la sala", userID))
	}
	room[userID] = c
	return nil
}

// Unregister elimina un cliente de la sala; cierra la conexión y elimina la sala si queda vacía.
// Además notifica a los demás usuarios (excepto al que sale).
func (h *Hub) Unregister(roomID, userID string) {
	h.mu.Lock()
	if room, ok := h.rooms[roomID]; ok {
		// Notificar salida a otros clientes
		for uid, client := range room {
			if uid != userID {
				client.send <- []byte(fmt.Sprintf("%s ha abandonado la sala", userID))
			}
		}
		// Cerrar y eliminar cliente
		if client, exists := room[userID]; exists {
			client.conn.Close()
			close(client.send)
			delete(room, userID)
		}
		debeCerrar := len(room) == 0
		h.mu.Unlock()

		if debeCerrar {
			h.CloseRoom(roomID)
		}
	}
}

// Broadcast envía un mensaje a todos los clientes de una sala.
func (h *Hub) Broadcast(roomID string, message []byte) {
	h.mu.Lock()
	room, ok := h.rooms[roomID]
	h.mu.Unlock()
	if !ok {
		return
	}
	for _, client := range room {
		select {
		case client.send <- message:
		default:
		}
	}
}

// CloseRoom cierra todas las conexiones en una sala y elimina la sala.
func (h *Hub) CloseRoom(roomID string) {
	h.mu.Lock()
	room, ok := h.rooms[roomID]
	if ok {
		delete(h.rooms, roomID)
	}
	h.mu.Unlock()

	for _, client := range room {
		client.send <- []byte(fmt.Sprintf("La sala %s ha sido cerrada por el servidor", roomID))
		client.conn.Close()
		close(client.send)
	}
}
