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

var ResponsesModel = Responses{}

func (r *Responses) Success(ctx *fiber.Ctx, data *fiber.Map) error {
	r.Status = 200
	r.Message = "success"
	r.Data = data
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) Error(ctx *fiber.Ctx, status int, message string) error {
	r.Status = status
	r.Message = message
	r.Data = nil
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) NotFound(ctx *fiber.Ctx) error {
	r.Status = 404
	r.Message = "not found"
	r.Data = nil
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) BadRequest(ctx *fiber.Ctx, message string) error {
	r.Status = 400
	r.Message = message
	r.Data = nil
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) Unauthorized(ctx *fiber.Ctx, message string) error {
	r.Status = 401
	r.Message = "unauthorized"
	r.Data = &fiber.Map{"message": message}
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) Forbidden(ctx *fiber.Ctx) error {
	r.Status = 403
	r.Message = "forbidden"
	r.Data = nil
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) InternalServerError(ctx *fiber.Ctx) error {
	r.Status = 500
	r.Message = "internal server error"
	r.Data = nil
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) ServiceUnavailable(ctx *fiber.Ctx) error {
	r.Status = 503
	r.Message = "service unavailable"
	r.Data = nil
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) GatewayTimeout(ctx *fiber.Ctx) error {
	r.Status = 504
	r.Message = "gateway timeout"
	r.Data = nil
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}

func (r *Responses) Custom(ctx *fiber.Ctx, status int, message string, data *fiber.Map) error {
	r.Status = status
	r.Message = message
	r.Data = data
	r.Timestamp = time.Now()
	return ctx.Status(r.Status).JSON(r)
}
