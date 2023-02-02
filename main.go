package main

import (
	"github.com/beesbuddy/beesbuddy-worker/cmd"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
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
	// TODO: Get name from env to have possibility to specify with different environment in mind
	config := core.NewConfig("dev")
	// inverse of control magic in context happens
	ctx := core.NewContext(config)

	(&cli.App{
		Name:    "BeesBuddy",
		Version: "main",
		Groups: []*cli.Group{
			{
				Name:  "make",
				Usage: "Make commands",
				Commands: []*cli.Command{
					{
						Name:       "migration",
						Usage:      "Run database migration",
						HandleFunc: cmd.Migrate(ctx),
					},
					{
						Name:  "user",
						Usage: "Make a user",
						Arguments: []string{
							"username",
							"email",
							"password",
						},
						HandleFunc: cmd.User(ctx),
					},
				},
			},
			{
				Name:  "web",
				Usage: "Web commands",
				Commands: []*cli.Command{
					{
						Name:       "serve",
						Usage:      "Serve the app",
						HandleFunc: cmd.WebServe(ctx),
					},
				},
			},
		},
	}).Execute()
}
