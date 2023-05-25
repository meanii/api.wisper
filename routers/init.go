package routers

import "github.com/gofiber/fiber/v2"

type Router struct {
	app *fiber.App
}

func (r *Router) Load(app *fiber.App) {
	r.app = app
	r.User() // loading up user router
}

func (r *Router) User() {
	userApp := r.app.Group("/user")
	user := User{}
	user.Init(userApp)
}
