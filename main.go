package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/meanii/api.wisper/clients"
	"github.com/meanii/api.wisper/configs"
	"github.com/meanii/api.wisper/routers"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: false,
		ServerHeader:  "wisper",
		AppName:       "wisper V0.0.1-dev.beta",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 status code
			code := fiber.StatusInternalServerError
			// Retreive the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			// Send custom error page
			return c.Status(code).JSON(fiber.Map{
				"message":   err.Error(),
				"status":    "error",
				"code":      code,
				"timestamp": time.Now(),
			})
		},
	})

	configs.InitConfig()

	// connecting to mongodb
	clients.MongoInit()

	// setting up middlewares
	app.Use(cors.New())

	// setting up request id
	app.Use(requestid.New())

	// setting up error handler
	app.Use(recover.New())

	// setting up logger
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${url}\n",
	}))

	// setting up limiter
	app.Use(limiter.New(limiter.Config{
		Max:                    20,
		Expiration:             30 * time.Second,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
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
			return (strings.Contains(c.Route().Path, "/ws")) || (c.Query("refresh") == "true")
		},
		CacheHeader: "x-wisper-cache",
		// Storage:      clients.RedisClient.Storage,
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

	PORT := fmt.Sprintf(":%s", configs.Env.Port)
	err := app.Listen(PORT)
	if err != nil {
		log.Panicf("could not start server, ERROR: %v", err)
	}
}
