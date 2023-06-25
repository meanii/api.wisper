package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/meanii/api.wisper/models"
	"github.com/meanii/api.wisper/utils"
)

type Auth struct {
	model models.User
	app   fiber.Router
}

func (a *Auth) Init(app fiber.Router) {
	a.app = app
	a.authorization()
	a.refresh()
}

func (a *Auth) authorization() {
	a.app.Get("/:token", func(c *fiber.Ctx) error {
		// get the user query
		token := c.Params("token")
		accessJwt := utils.JWT[utils.AccessTokenRawPayload]{}

		validateToken, err := accessJwt.ValidateToken(token)
		if err != nil {
			return utils.ResponsesModel.Unauthorized(c, err.Error())
		}
		return utils.ResponsesModel.Success(c, &fiber.Map{
			"id":       validateToken.Payload.ID,
			"username": validateToken.Payload.Username,
		})
	})
}

func (a *Auth) refresh() {
	a.app.Post("/refresh", func(c *fiber.Ctx) error {
		token := &utils.Tokens{}
		if err := c.BodyParser(token); err != nil {
			return utils.ResponsesModel.Error(c, fiber.StatusBadRequest, err.Error())
		}

		tokens, err := utils.RefreshTokens(*token)
		if err != nil {
			return utils.ResponsesModel.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.ResponsesModel.Success(c, &fiber.Map{"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
	})
}
