package app

import (
	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app/settings"
	"github.com/leonidasdeim/goconfig"
	"github.com/leonidasdeim/goconfig/pkg/filehandler"
)

func NewConfig(configName string) *goconfig.Config[settings.AppSettings] {
	h, _ := filehandler.New(filehandler.WithName(configName))
	cfg, err := goconfig.Init[settings.AppSettings](h)

	if err != nil {
		panic("Unable to load config")
	}

	cfg.AddSubscriber(internal.WorkerKey)

	return cfg
}
