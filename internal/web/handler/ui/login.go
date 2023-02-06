package ui

import (
	"net/http"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/petaki/support-go/forms"
)

func WebGetLoginHandler(ctx *app.Ctx) fiber.Handler {
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ctx.InertiaManager.Render(w, r, "auth/Login", map[string]interface{}{
			"errors": forms.Bag{},
		})

		if err != nil {
			serverError(ctx, w, err)
		}
	})
}

func WebPostLoginHandler(ctx *app.Ctx) fiber.Handler {
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		form, err := forms.NewFromRequest(w, r)
		if err != nil {
			formError(ctx, w, err)

			return
		}

		form.Required("username_or_email", "password")

		if form.IsValid() {
			user, err := ctx.UserRepository.Authenticate(form.Data["username_or_email"].(string), form.Data["password"].(string))

			if err != nil {
				serverError(ctx, w, err)

				return
			}

			if err != nil {
				serverError(ctx, w, err)

				return
			}

			if form.Data["remember"].(bool) {
				if user.RememberToken == "" {
					token, err := internal.GenerateToken()
					if err != nil {
						serverError(ctx, w, err)

						return
					}

					err = ctx.UserRepository.UpdateRememberToken(user, token)
					if err != nil {
						serverError(ctx, w, err)

						return
					}
				}

				err = ctx.RememberCookie.SetValue(w, user.RememberCookie())
				if err != nil {
					serverError(ctx, w, err)

					return
				}
			}
		}

		err = ctx.InertiaManager.Render(w, r, "auth/Login", map[string]interface{}{
			"errors": form.Errors,
		})

		if err != nil {
			serverError(ctx, w, err)
		}

	})
}
