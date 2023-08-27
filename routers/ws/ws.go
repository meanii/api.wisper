package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app fiber.Router
}

func (w *Server) Init(app fiber.Router) {
	w.app = app
	// initialize root ws server here
	// ws://localhost:3000/wisper/ws
	w.Server()

	// initialize echo ws server here
	// ws://localhost:3000/wisper/ws/echo
	w.Echo()
}

func (w *Server) Server() {
	w.app.Use("/", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}
