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
			return utils.ResponsesModel.Error(c, http.StatusBadRequest, fmt.Sprintf("something went wrong! ERROR: %s", err.Error()))
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			return utils.ResponsesModel.Error(c, http.StatusBadRequest, fmt.Sprintf("something went wrong! ERROR: %s", validationErr.Error()))
		}

		// check if user already exists
		var userExists models.User
		if err := models.UserModel.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userExists); err == nil {
			return utils.ResponsesModel.Error(c, http.StatusBadRequest, fmt.Sprintf("user already exists!"))
		}

		// encrypt password
		hash, err := utils.Hash(user.Password)
		if err != nil {
			return utils.ResponsesModel.Error(c, http.StatusInternalServerError, fmt.Sprintf("something went wrong! ERROR: %s", err.Error()))
		}
		user.Password = hash

		// insert user
		result, err := models.UserModel.InsertOne(ctx, user)

		if err != nil {
			return utils.ResponsesModel.Error(c, http.StatusInternalServerError, fmt.Sprintf("something went wrong! ERROR: %s", err.Error()))
		}
		return utils.ResponsesModel.Success(c, &fiber.Map{"id": result.InsertedID})
	})
}

func (u *User) login() {
	u.app.Post("/login", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user = u.model
		defer cancel()

		if err := c.BodyParser(&user); err != nil {
			return utils.ResponsesModel.Error(c, http.StatusBadRequest, fmt.Sprintf("something went wrong! ERROR: %s", err.Error()))
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			return utils.ResponsesModel.Error(c, http.StatusBadRequest, "validation has been failed! ERROR: "+validationErr.Error())
		}

		result := models.UserModel.FindOne(ctx, bson.M{"username": user.Username})
		if result.Err() != nil {
			return utils.ResponsesModel.Error(c, http.StatusUnauthorized, "username or password is incorrect!")
		}

		var foundUser models.User
		err := result.Decode(&foundUser)

		if err != nil {
			return utils.ResponsesModel.Error(c, http.StatusInternalServerError, fmt.Sprintf("something went wrong! ERROR: %s", err.Error()))
		}

		// validate password
		if err := utils.Verify(user.Password, foundUser.Password); err != nil {
			return utils.ResponsesModel.Error(c, http.StatusUnauthorized, "username or password is incorrect!")
		}

		// generate token
		accessJwt := utils.JWT[utils.AccessTokenRawPayload]{}
		refreshJwt := utils.JWT[utils.RefreshTokenRawPayload]{}

		// generate access token payload
		accessTokenPayload := &utils.AccessTokenRawPayload{
			ID:       foundUser.Id,
			Username: user.Username,
		}

		// generate access token
		accessToken, err := accessJwt.GenerateToken(*accessTokenPayload, 1)
		if err != nil {
			return utils.ResponsesModel.Error(c, http.StatusInternalServerError, fmt.Sprintf("something went wrong! ERROR: %s", err.Error()))
		}

		// generate refresh token payload
		refreshTokenPayload := &utils.RefreshTokenRawPayload{
			AccessToken: accessToken,
		}

		// generate refresh token
		refreshToken, err := refreshJwt.GenerateToken(*refreshTokenPayload, 24)
		if err != nil {
			return utils.ResponsesModel.Error(c, http.StatusInternalServerError, fmt.Sprintf("something went wrong! ERROR: %s", err.Error()))
		}

		return utils.ResponsesModel.Success(c, &fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken})
	})
}
