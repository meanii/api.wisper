package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/meanii/api.wisper/configs"
	"github.com/meanii/api.wisper/utils"
)

type Status struct {
	app fiber.Router
}

func (a *Status) Init(app fiber.Router) {
	a.app = app
	a.status()
}

func (a *Status) status() {
	a.app.Get("/", func(c *fiber.Ctx) error {

		return utils.ResponsesModel.Success(c, &fiber.Map{
			"status":  "ok",
			"version": configs.WisperVersion,
			"message": "wisper is running",
		})
	})
}
