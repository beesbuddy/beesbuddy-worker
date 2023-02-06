package util

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTToken(appClientKey, secret string) (string, error) {

	claims := jwt.MapClaims{
		"app_key": appClientKey,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
