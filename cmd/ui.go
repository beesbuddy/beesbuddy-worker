package cmd

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type uiModule struct {
	app     *fiber.App
	embedFS *embed.FS
}

func NewUIRunner(app *fiber.App, embedFS *embed.FS) core.ModuleRunner {
	module := &uiModule{app: app, embedFS: embedFS}

	return module
}

func (mod *uiModule) Run() {
	publicFS, err := fs.Sub(mod.embedFS, "public")

	if err != nil {
		log.Fatal(err)
	}

	group := mod.app.Group("", cors.New(), csrf.New())

	group.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(publicFS),
		NotFoundFile: "index.html",
	}))
}

func (mod *uiModule) CleanUp() {

}
