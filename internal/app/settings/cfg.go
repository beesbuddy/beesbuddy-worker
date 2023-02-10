package settings

import (
	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/leonidasdeim/goconfig"
	"github.com/leonidasdeim/goconfig/pkg/filehandler"
)

func NewConfig(configName string) *goconfig.Config[AppSettings] {
	h, _ := filehandler.New(filehandler.WithName(configName))
	cfg, err := goconfig.Init[AppSettings](h)

	if err != nil {
		panic("Unable to load config")
	}

	cfg.AddSubscriber(internal.WorkerKey)

	return cfg
}
