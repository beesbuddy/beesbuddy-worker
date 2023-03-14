package main

import (
	"github.com/beesbuddy/beesbuddy-worker/app"
	"github.com/beesbuddy/beesbuddy-worker/cmd"
	"github.com/beesbuddy/beesbuddy-worker/constants"
	"github.com/beesbuddy/beesbuddy-worker/pref"
	"github.com/beesbuddy/beesbuddy-worker/util"
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
	// Initialize configuration with default values
	configPath := util.GetEnv("BB_CONFIG_PATH", "./")
	environmentName := util.GetEnv("BB_ENV_NAME", "dev")
	persistedDataPath := util.GetEnv("BB_DATA_PATH", "./data")

	config := pref.NewPreferences[pref.AppPreferences](
		configPath,
		environmentName,
	)

	// Update configuration with env variables
	newConfig := config.GetConfig()
	newConfig.StoragePath = persistedDataPath
	config.UpdateConfig(newConfig)

	// Create applicaiton context, inverse of control magic in context happens here
	appCtx := app.NewContext(config)

	// Add subscriber for worker
	config.AddSubscriber(constants.WorkerKey)

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
