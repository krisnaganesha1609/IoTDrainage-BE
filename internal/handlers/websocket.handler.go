package handlers

import (
	"github.com/gofiber/contrib/v3/websocket"
)

func (h *Handler) ConnectWebsocket(c *websocket.Conn) {
	h.Service.Register(c)
	defer func() {
		h.Service.Unregister(c)
		c.Close()
	}()
	for {
		// keeping the connection alive by reading messages, but we can ignore the content
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
