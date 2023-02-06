package core

import (
	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	c "github.com/leonidasdeim/goconfig"
	fh "github.com/leonidasdeim/goconfig/pkg/filehandler"
)

func NewConfig(configName string) *c.Config[model.Config] {
	h, _ := fh.New(fh.WithName(configName))
	cfg, err := c.Init[model.Config](h)

	if err != nil {
		panic("Unable to load config")
	}

	cfg.AddSubscriber(internal.WorkerKey)

	return cfg
}
