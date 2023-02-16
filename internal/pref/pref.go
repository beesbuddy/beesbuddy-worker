package pref

import "time"

type AppConfig struct {
	AppName             string       `default:"BeesBuddy Worker"`
	AppHost             string       `default:"0.0.0.0"`
	AppPort             int          `default:"4000"`
	Clients             []Client     `required:"true"`
	Secret              string       `required:"true"`
	IsPrefork           bool         `default:"false"`
	IsProd              bool         `default:"false"`
	BrokerTCPUrl        string       `requred:"true"`
	Subscribers         []Subscriber `required:"false"`
	StoragePath         string       `default:"./data"`
	InfluxDbAccessToken string       `default:"change_it"`
	InfluxDbURL         string       `default:"http://localhost:8086"`
	InfluxDbOrg         string       `default:"BeesBuddy"`
}

type Client struct {
	AppKey string `required:"true"`
}

type Subscriber struct {
	ApiaryId  string    `required:"true"`
	HiveId    string    `required:"true"`
	CreatedAt time.Time `required:"true"`
}

type Preferences[T any] interface {
	UpdateConfig(T)
	GetConfig() T
	AddSubscriber(string)
	GetSubscriber(string) <-chan bool
}