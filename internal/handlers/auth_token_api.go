package handlers

import (
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
)

func getClient(appKey string) (models.Client, bool) {
	return lo.Find(core.GetCfg().Clients, func(c models.Client) bool {
		return c.AppKey == appKey
	})
}

// Generate api key (token) by app key and secret key
// @Summary Authenticate client
// @Description Create a token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} core.ResponseHTTP{data=string}
// @Failure 503 {object} core.ResponseHTTP{}
// @Param dto.ClientInput body dto.ClientInput true "ClientInput"
// @Router /auth/token [post]
func ApiGenerateToken(ctx *fiber.Ctx) error {
	var input dto.ClientInput

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on token generate request", "data": err})
	}

	appKey := input.AppKey

	client, ok := getClient(appKey)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	claims := jwt.MapClaims{
		"app_key": client.AppKey,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(core.GetCfg().Secret))

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(core.ResponseHTTP{Success: true, Message: "Successfully generated token", Data: t})
}
