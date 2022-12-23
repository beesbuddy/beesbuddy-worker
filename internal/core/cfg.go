package core

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	config "github.com/leonidasdeim/goconfig"
)

var configModel *model.Config
var appConfig *config.Config[model.Config]

func GetCfgModel() *model.Config {
	return configModel
}

func GetAppCfgObject() *config.Config[model.Config] {
	return appConfig
}

func InitializeConfig() {
	appConfig, err := config.Init[model.Config](config.WithName("dev"))

	if err != nil {
		panic("Unable to load config")
	}

	configModel = appConfig.GetCfg()
}
