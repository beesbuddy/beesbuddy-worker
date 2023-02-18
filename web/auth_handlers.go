package web

import (
	"github.com/beesbuddy/beesbuddy-worker/dto"
	"github.com/gofiber/fiber/v2"
)

func AuthError(f *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return f.Status(fiber.StatusBadRequest).
			JSON(&dto.ResponseHTTP{
				Success: false,
				Data:    nil,
				Message: "Missing or malformed token",
			})
	}

	return f.Status(fiber.StatusUnauthorized).
		JSON(&dto.ResponseHTTP{
			Success: false,
			Data:    nil,
			Message: "Invalid or expired token",
		})
}
