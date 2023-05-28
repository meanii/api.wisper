package routers

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/meanii/api.wisper/models"
	"github.com/meanii/api.wisper/utils"
	"go.mongodb.org/mongo-driver/bson"
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
	u.login()
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

		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.Responses{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("something went wrong! ERROR: %s", err.Error()),
			})
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.Responses{
				Status:  http.StatusBadRequest,
				Message: "validation has been failed! ERROR: " + validationErr.Error(),
			})
		}

		// check if user already exists
		var userExists models.User
		if err := models.UserModel.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userExists); err == nil {
			return c.Status(http.StatusBadRequest).JSON(utils.Responses{
				Status:  http.StatusBadRequest,
				Message: "user already exists!",
			})
		}

		// encrypt password
		hash, err := utils.Hash(user.Password)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(utils.Responses{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("something went wrong! ERROR: %s", err.Error()),
			})
		}
		user.Password = hash

		// insert user
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
			Data:    &fiber.Map{"id": result.InsertedID},
		})

	})
}

func (u *User) login() {
	u.app.Post("/login", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user = u.model
		defer cancel()

		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.Responses{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("something went wrong! ERROR: %s", err.Error()),
			})
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.Responses{
				Status:  http.StatusBadRequest,
				Message: "validation has been failed! ERROR: " + validationErr.Error(),
			})
		}

		result := models.UserModel.FindOne(ctx, bson.M{"username": user.Username})
		if result.Err() != nil {
			return c.Status(http.StatusUnauthorized).JSON(utils.Responses{
				Status:  http.StatusUnauthorized,
				Message: "username or password is incorrect!",
			})
		}

		var foundUser models.User
		err := result.Decode(&foundUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(utils.Responses{
				Status:  http.StatusInternalServerError,
				Message: "something went wrong!",
			})
		}

		// validate password
		if err := utils.Verify(user.Password, foundUser.Password); err != nil {
			return c.Status(http.StatusUnauthorized).JSON(utils.Responses{
				Status:  http.StatusUnauthorized,
				Message: "username or password is incorrect!",
			})
		}

		// generate token
		var jwt = utils.JWT{}
		// generate payload
		var payload = &utils.Payload{
			ID:       user.Id,
			Username: user.Username,
		}

		// generate access token
		accessToken, err := jwt.GenerateAccessToken(payload, 2)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(utils.Responses{
				Status:  http.StatusInternalServerError,
				Message: "something went wrong!",
			})
		}

		// generate refresh token
		refreshToken, err := jwt.GenerateRefreshToken(payload, 2)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(utils.Responses{
				Status:  http.StatusInternalServerError,
				Message: "something went wrong!",
			})
		}

		return c.Status(http.StatusOK).JSON(utils.Responses{
			Status:  http.StatusOK,
			Message: "user has been logged in!",
			Data:    &fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken},
		})

	})
}
