package core

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	c "github.com/leonidasdeim/goconfig"
)

var cfgObject *c.Config[model.Config]

func GetCfgModel() model.Config {
	return cfgObject.GetCfg()
}

func GetCfgObject() *c.Config[model.Config] {
	return cfgObject
}

func InitializeCfg() {
	cfg, err := c.Init[model.Config](c.WithName("dev"))

	if err != nil {
		panic("Unable to load config")
	}

	cfgObject = cfg
}
