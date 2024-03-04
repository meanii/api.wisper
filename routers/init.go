package routers

import (
	"github.com/gofiber/fiber/v2"

	wws "github.com/meanii/api.wisper/routers/ws"
)

type Router struct {
	app  *fiber.App
	root fiber.Router
}

func (r *Router) Load(app *fiber.App) {
	r.app = app
	r.root = r.app.Group("/wisper")
	r.rootWelcome() // loading up /wisper welcome message
	r.Status()      // loading up status router
	r.User()        // loading up user router
	r.Auth()        // loading up auth router
	r.WS()          // loading up ws router
}

func (r *Router) rootWelcome() {
	r.root.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("welcome to wisper api ☂️.")
	})
}

func (r *Router) WS() {
	wsApp := r.root.Group("/ws")
	ws := wws.Server{}
	ws.Init(wsApp)
}

func (r *Router) User() {
	userApp := r.root.Group("/user")
	user := User{}
	user.Init(userApp)
}

func (r *Router) Auth() {
	authApp := r.root.Group("/auth")
	auth := Auth{}
	auth.Init(authApp)
}

func (r *Router) Status() {
	statusApp := r.root.Group("/status")
	status := Status{}
	status.Init(statusApp)
}
