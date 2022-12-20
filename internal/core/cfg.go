package core

import (
	config "github.com/beesbuddy/beesbuddy-config"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
)

var appConfig *model.Config

func GetConfig() *model.Config {
	return appConfig
}

func InitializeConfig() {
	initializedConfig, error := config.Init[model.Config](config.WithName("dev"))

	if error != nil {
		panic("Unable to load config")
	}

	appConfig = initializedConfig.GetCfg()
}
