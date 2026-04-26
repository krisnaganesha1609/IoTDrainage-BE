package services

import (
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func (s *Service) Register(c *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	clients[c] = true
}

func (s *Service) Unregister(c *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, c)
}

func (s *Service) Broadcast(data interface{}) {
	mu.Lock()
	defer mu.Unlock()

	for c := range clients {
		c.WriteJSON(data)
	}
}
