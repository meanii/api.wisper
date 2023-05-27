package routers

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/meanii/api.wisper/models"
	"github.com/meanii/api.wisper/utils"
	"net/http"
	"time"
)

var validate = validator.New()

type User struct {
	model models.User
	app   fiber.Router
}

func (u *User) Init(app fiber.Router) {
	u.app = app
	u.welcome()
	u.signup()
}

func (u *User) welcome() {
	u.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("welcome to user home.")
	})
}

func (u *User) signup() {

	u.app.Post("/signup", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user = u.model
		defer cancel()

		if validationErr := validate.Struct(user); validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.Responses{
				Status:  http.StatusBadRequest,
				Message: "validation has been failed!",
			})
		}

		if err := c.BodyParser(&user); err != nil {
			return err
		}

		result, err := models.UserModel.InsertOne(ctx, user)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(utils.Responses{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("something went wrong! ERROR: %s", err.Error()),
			})
		}

		return c.Status(201).JSON(utils.Responses{
			Status:  http.StatusCreated,
			Message: "user has been created!",
			Data:    &fiber.Map{"data": result},
		})

	})
}

func (u *User) login() {
	u.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("")
	})
}
