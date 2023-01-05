package handlers

import (
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func getUserByUsername(u string) (*models.User, error) {
	user := &models.User{
		IsEnabled: true,
	}

	return user, nil
}

// Login is a function to login
// @Summary Authenticate user
// @Description Create a user
// @Tags login
// @Accept json
// @Produce json
// @Success 200 {object} core.ResponseHTTP{data=string}
// @Failure 503 {object} core.ResponseHTTP{}
// @Param dto.UserLoginInput body dto.UserLoginInput true "UserOutput"
// @Router /auth/login [post]
func ApiLogin() func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		var input dto.UserLoginInput

		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
		}
		identity := input.Username
		password := input.Password

		user, err := getUserByUsername(identity)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Error on username", "data": err})
		}

		if user == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "UserOutput not found", "data": err})
		}

		if !c.CheckPasswordHash(password, user.Password) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(core.ResponseHTTP{
				Success: false,
				Message: "Invalid password",
			})
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["user_id"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte(c.GetCfg().Secret))
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.JSON(core.ResponseHTTP{Success: true, Message: "Success login", Data: t})
	}
}
