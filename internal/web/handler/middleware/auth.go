package middleware

import (
	"context"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	"github.com/gofiber/fiber/v2"
)

func RedirectIfNotAuthenticated(ctx *app.Ctx) func(*fiber.Ctx) error {
	return func(f *fiber.Ctx) error {
		session, err := ctx.SessionManager.Get(f)

		if err != nil {
			panic(err)
		}

		if authUser(f.UserContext()) == nil {
			session.Set(internal.SessionIntendedURLKey, f.Request().URI().Path())
			return f.Redirect("login")
		}

		return f.Next()
	}
}

func authUser(ctx context.Context) *model.User {
	user, ok := ctx.Value(internal.AuthUserKey).(*model.User)
	if !ok {
		return nil
	}

	return user
}
