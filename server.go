package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/meanii/api.wisper/clients"
	"github.com/meanii/api.wisper/routers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	configs "github.com/meanii/api.wisper/configs"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: false,
		ServerHeader:  "Wisper",
		AppName:       "Wisper V0.0.1-dev.beta",
	})

	// load the .env file
	env := configs.GetConfig()

	// connecting to mongodb
	_ = clients.GetClient()

	// setting up middlewares
	app.Use(cors.New())
	// setting up error handler
	app.Use(recover.New())

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	router := routers.Router{}
	router.Load(app)

	// root router
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("welcome to wisper api gatway ☂️.")
	})

	PORT := fmt.Sprintf(":%s", env.Port)
	err := app.Listen(PORT)
	if err != nil {
		log.Panic("could not start server, ERROR: ", err)
	}

}
