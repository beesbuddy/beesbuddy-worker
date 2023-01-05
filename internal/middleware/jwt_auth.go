package middleware

import (
	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v3"
)

// JwtProtected protect routes
func JwtProtected() fiber.Handler {
	return jwtMiddleware.New(jwtMiddleware.Config{
		SigningKey:   []byte(c.GetCfg().Secret),
		ErrorHandler: jwtError,
	})
}

func jwtError(ctx *fiber.Ctx, err error) error {

	if err.Error() == "Missing or malformed JWT" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(&c.ResponseHTTP{
				Success: false,
				Data:    nil,
				Message: "Missing or malformed JWT",
			})
	}
	return ctx.Status(fiber.StatusUnauthorized).
		JSON(&c.ResponseHTTP{
			Success: false,
			Data:    nil,
			Message: "Invalid or expired JWT",
		})
}
