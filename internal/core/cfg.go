package core

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	c "github.com/leonidasdeim/goconfig"
)

var cfgObject *c.Config[models.Config]

func GetCfg() models.Config {
	return cfgObject.GetCfg()
}

func GetCfgObject() *c.Config[models.Config] {
	return cfgObject
}

func init() {
	cfg, err := c.Init[models.Config](c.WithName("dev"))

	if err != nil {
		panic("Unable to load config")
	}

	cfg.AddSubscriber(WorkerKey)

	cfgObject = cfg
}
