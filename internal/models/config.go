package models

import "github.com/beesbuddy/beesbuddy-worker/internal/dto"

type Config struct {
	AppName        string        `default:"BeesBuddy Worker"`
	AppHost        string        `default:"0.0.0.0"`
	AppPort        int           `default:"4000"`
	Clients        []Client      `required:"true"`
	Secret         string        `required:"true"`
	UiHotReloadUrl string        `default:"http://localhost:8080"`
	IsPrefork      bool          `default:"false"`
	IsProd         bool          `default:"false"`
	BrokerTCPUrl   string        `requred:"true"`
	Admin          dto.UserInput `required:"true"`
	Subscribers    []Subscriber  `required:"false"`
}

type Client struct {
	AppKey string `required:"true"`
}
