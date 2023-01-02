package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	handler "github.com/beesbuddy/beesbuddy-worker/internal/handler/web"
	"github.com/gofiber/adaptor/v2"
)

type webCmd struct {
	app *core.App
}

func NewWebRunner(app *core.App) core.CmdRunner {
	module := &webCmd{app: app}
	return module
}

func (cmd *webCmd) Run() {
	cfg := core.GetCfg()

	web := cmd.app.Router.Group("/")
	web.Use(adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return cmd.app.InertiaManager.Middleware(next)
	}))

	web.Get("/", handler.HomeHandler(cmd.app))

	go func() {
		if err := cmd.app.Router.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Panic(err)
		}
	}()
}

func (cmd *webCmd) CleanUp() {
	log.Println("Gracefully closing web...")
	go func() {
		err := cmd.app.Router.Shutdown()

		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
	}()
}
