//go:build mocking
// +build mocking

package mocking

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockWSClient struct {
	mu       sync.Mutex
	messages [][]byte
	closed   bool
}

func (c *MockWSClient) Send(msg []byte) {
	c.mu.Lock()
	c.messages = append(c.messages, msg)
	c.mu.Unlock()
}

func (c *MockWSClient) Close() {
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()
}

// HubAdapter replica la lógica del hub con mocks simples.
type HubAdapter struct {
	mu         sync.Mutex
	rooms      map[string]map[string]*MockWSClient
	maxClients int
}

func NewHubAdapter(max int) *HubAdapter {
	return &HubAdapter{rooms: make(map[string]map[string]*MockWSClient), maxClients: max}
}

func (h *HubAdapter) Register(roomID, userID string, c *MockWSClient) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[string]*MockWSClient)
	}
	if len(h.rooms[roomID]) >= h.maxClients {
		return assert.AnError
	}
	for _, client := range h.rooms[roomID] {
		client.Send([]byte(userID + " joined"))
	}
	h.rooms[roomID][userID] = c
	return nil
}

func (h *HubAdapter) Unregister(roomID, userID string) {
	h.mu.Lock()
	room, ok := h.rooms[roomID]
	if !ok {
		h.mu.Unlock()
		return
	}
	for uid, client := range room {
		if uid != userID {
			client.Send([]byte(userID + " left"))
		}
	}
	target := room[userID]
	delete(room, userID)
	empty := len(room) == 0
	h.mu.Unlock()

	if target != nil {
		target.Close()
	}
	if empty {
		h.CloseRoom(roomID)
	}
}

func (h *HubAdapter) Broadcast(roomID string, msg []byte) {
	h.mu.Lock()
	room := h.rooms[roomID]
	h.mu.Unlock()
	for _, client := range room {
		client.Send(msg)
	}
}

func (h *HubAdapter) CloseRoom(roomID string) {
	h.mu.Lock()
	room := h.rooms[roomID]
	delete(h.rooms, roomID)
	h.mu.Unlock()

	for _, client := range room {
		client.Send([]byte("closed"))
		client.Close()
	}
}

func TestHubAdapter_RegisterAndBroadcast(t *testing.T) {
	hub := NewHubAdapter(2)
	c1 := &MockWSClient{}
	c2 := &MockWSClient{}

	err := hub.Register("room", "u1", c1)
	assert.NoError(t, err)

	err = hub.Register("room", "u2", c2)
	assert.NoError(t, err)
	assert.Len(t, c1.messages, 1) // u2 joined notification

	hub.Broadcast("room", []byte("hello"))
	assert.Equal(t, []byte("hello"), c1.messages[len(c1.messages)-1])
	assert.Equal(t, []byte("hello"), c2.messages[len(c2.messages)-1])
}

func TestHubAdapter_Register_Full(t *testing.T) {
	hub := NewHubAdapter(1)
	err := hub.Register("room", "u1", &MockWSClient{})
	assert.NoError(t, err)

	err = hub.Register("room", "u2", &MockWSClient{})
	assert.Error(t, err)
}

func TestHubAdapter_UnregisterAndClose(t *testing.T) {
	hub := NewHubAdapter(3)
	c1 := &MockWSClient{}
	c2 := &MockWSClient{}

	_ = hub.Register("room", "u1", c1)
	_ = hub.Register("room", "u2", c2)

	hub.Unregister("room", "u1")
	assert.True(t, c1.closed)
	assert.Len(t, c2.messages, 1)
}
