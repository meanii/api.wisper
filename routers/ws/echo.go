package ws

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
)

func (w *Server) Echo() {
	w.app.Get("/echo", websocket.New(func(c *websocket.Conn) {
		// websocket.Conn bindings:
		var (
			mt  int
			msg []byte
			err error
		)

		// Read message from client
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break
			}

			// Print the message to the console
			fmt.Println("ws://localhost:3000/wisper/ws/echo", string(msg))

			// Write message back to client
			if err = c.WriteMessage(mt, msg); err != nil {
				break
			}
		}
	}))
}
