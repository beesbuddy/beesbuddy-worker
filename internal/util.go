package internal

import (
	"encoding/hex"
	"os"
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	"github.com/chmike/securecookie"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

func GenerateToken() (string, error) {
	bytes, err := securecookie.GenerateRandomKey()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func GenerateJWTToken(appClientKey, secret string) (string, error) {

	claims := jwt.MapClaims{
		"app_key": appClientKey,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
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

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
