package handlers

import (
	"net/http"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func HomeHandler(app *core.App) fiber.Handler {
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := app.InertiaManager.Render(w, r, "home/Index", map[string]interface{}{
			"appName": core.GetCfg().AppName,
		})

		if err != nil {
			panic("unable to render home page")
		}
	})
}
