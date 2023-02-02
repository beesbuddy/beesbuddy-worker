package core

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	c "github.com/leonidasdeim/goconfig"
	fh "github.com/leonidasdeim/goconfig/pkg/filehandler"
)

var cfgObject *c.Config[models.Config]

func GetCfg() models.Config {
	return cfgObject.GetCfg()
}

func GetCfgObject() *c.Config[models.Config] {
	return cfgObject
}

func init() {
	h, _ := fh.New(fh.WithName("dev"))
	cfg, err := c.Init[models.Config](h)

	if err != nil {
		panic("Unable to load config")
	}

	cfg.AddSubscriber(WorkerKey)

	cfgObject = cfg
}
