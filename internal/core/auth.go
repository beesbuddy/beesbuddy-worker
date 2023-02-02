package core

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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
