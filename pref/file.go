package pref

import (
	"github.com/beesbuddy/beesbuddy-worker/constants"
	"github.com/beesbuddy/beesbuddy-worker/internal/log"
	"github.com/leonidasdeim/goconfig"
	"github.com/leonidasdeim/goconfig/pkg/filehandler"
)

type FilePreferences[T any] struct {
	fileConfig *goconfig.Config[T]
}

func NewPreferences[T any](configName string) *FilePreferences[T] {
	h, _ := filehandler.New(filehandler.WithName(configName))
	cfg, err := goconfig.Init[T](h)

	if err != nil {
		log.Error.Fatal("unable to load config")
		panic("unable to load config")
	}

	// Add default subscriber
	cfg.AddSubscriber(constants.WorkerKey)

	return &FilePreferences[T]{
		fileConfig: cfg,
	}
}

func (f *FilePreferences[T]) AddSubscriber(key string) {
	f.fileConfig.AddSubscriber(key)
}

func (f FilePreferences[T]) GetSubscriber(key string) <-chan bool {
	return f.fileConfig.GetSubscriber(key)
}

func (f FilePreferences[T]) GetConfig() T {
	return f.fileConfig.GetCfg()
}

func (f *FilePreferences[T]) UpdateConfig(config T) {
	f.fileConfig.Update(config)
}
