package cmd

import (
	"log"
	"net/http"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	handler "github.com/beesbuddy/beesbuddy-worker/internal/handler/web"
	"github.com/gofiber/adaptor/v2"
)

type webCmd struct {
	app *core.App
}

func NewWebRunner(app *core.App) core.CmdRunner {
	cmd := &webCmd{app: app}

	return cmd
}

func (cmd *webCmd) Run() {
	web := cmd.app.Router.Group("/")
	web.Use(adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return cmd.app.InertiaManager.Middleware(next)
	}))

	web.Get("/", handler.HomeHandler(cmd.app))
}

func (cmd *webCmd) CleanUp() {
	log.Println("Gracefully closing web...")
}
