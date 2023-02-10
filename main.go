package main

import (
	"github.com/beesbuddy/beesbuddy-worker/cmd"
	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/app/settings"
	"github.com/petaki/support-go/cli"
)

// @title BeesBuddy worker
// @version 1.0
// @description This is an API for Worker Module

// @contact.name Viktor Nareiko
// @contact.email vnareiko.lt@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config := settings.NewConfig(internal.GetEnv("BEESBUDDY_ENV", "dev"))
	// inverse of control magic in context happens
	appCtx := app.NewContext(config)

	(&cli.App{
		Name:    "BeesBuddy",
		Version: "main",
		Groups: []*cli.Group{
			{
				Name:  "make",
				Usage: "Make commands",
				Commands: []*cli.Command{
					{
						Name:  "token",
						Usage: "Make a jwt token",
						Arguments: []string{
							"key",
						},
						HandleFunc: cmd.Token(appCtx),
					},
				},
			},
			{
				Name:  "worker",
				Usage: "Worker commands",
				Commands: []*cli.Command{
					{
						Name:       "serve",
						Usage:      "Serve the app",
						HandleFunc: cmd.WrokerServe(appCtx),
					},
				},
			},
		},
	}).Execute()
}
