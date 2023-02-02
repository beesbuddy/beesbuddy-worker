package core

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	c "github.com/leonidasdeim/goconfig"
	fh "github.com/leonidasdeim/goconfig/pkg/filehandler"
)

func NewConfig(configName string) *c.Config[models.Config] {
	h, _ := fh.New(fh.WithName(configName))
	cfg, err := c.Init[models.Config](h)

	if err != nil {
		panic("Unable to load config")
	}

	cfg.AddSubscriber(WorkerKey)

	return cfg
}
