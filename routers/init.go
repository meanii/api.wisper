package routers

import "github.com/gofiber/fiber/v2"

type Router struct {
	app  *fiber.App
	root fiber.Router
}

func (r *Router) Load(app *fiber.App) {
	r.app = app
	r.root = r.app.Group("/wisper")
	r.rootWelcome() // loading up /wisper welcome message
	r.User()        // loading up user router
}

func (r *Router) rootWelcome() {
	r.root.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("welcome to wisper api ☂️.")
	})
}

func (r *Router) User() {
	userApp := r.root.Group("/user")
	user := User{}
	user.Init(userApp)
}
