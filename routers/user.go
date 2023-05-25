package routers

import "github.com/gofiber/fiber/v2"

type User struct {
	app fiber.Router
}

func (u *User) Init(app fiber.Router) {
	u.app = app
	u.welcome()
}

func (u *User) welcome() {
	u.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("welcome to user home.")
	})
}

func (u *User) signup() {
}
