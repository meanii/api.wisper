package utils

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type Responses struct {
	Status    int        `json:"status"`
	Message   string     `json:"message"`
	Data      *fiber.Map `json:"data"`
	Timestamp time.Time  `json:"timestamp"`
}
