package api

import (
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
)

// Generate api key (token) by app key and secret key
// @Summary Authenticate client
// @Description Create a token
// @Tags auth
// @Accept json
// @Deprecated true
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=string}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Param dto.ClientInput body dto.ClientInput true "ClientInput"
// @Router /auth/token [post]
func ApiGenerateToken(ctx *app.Ctx) fiber.Handler {
	return func(f *fiber.Ctx) error {
		var input dto.ClientInput

		if err := f.BodyParser(&input); err != nil {
			return f.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on token generate request", "data": err})
		}

		appKey := input.AppKey

		client, ok := getClient(appKey, ctx.Config.GetCfg().Clients)

		if !ok {
			return f.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
		}

		claims := jwt.MapClaims{
			"app_key": client.AppKey,
			"exp":     time.Now().Add(time.Minute * 30).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(ctx.Config.GetCfg().Secret))

		if err != nil {
			return f.SendStatus(fiber.StatusInternalServerError)
		}

		return f.JSON(dto.ResponseHTTP{Success: true, Message: "Successfully generated token", Data: t})
	}
}

func getClient(appKey string, clients []model.Client) (model.Client, bool) {
	return lo.Find(clients, func(client model.Client) bool {
		return client.AppKey == appKey
	})
}
