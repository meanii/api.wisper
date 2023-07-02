package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/meanii/api.wisper/clients"
	"github.com/meanii/api.wisper/routers"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	configs "github.com/meanii/api.wisper/configs"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
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

	// setting up request id
	app.Use(requestid.New())

	// setting up error handler
	app.Use(recover.New())

	// setting up logger
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// setting up limiter
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests, please try again later.",
			})
		},
	}))

	// setting up cache
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	// setting up compress
	app.Use(compress.New())

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
