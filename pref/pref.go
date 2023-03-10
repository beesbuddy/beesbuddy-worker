package pref

import "time"

type AppPreferences struct {
	AppName             string       `default:"BeesBuddy Worker"`
	AppHost             string       `default:"lohalhost"`
	AppPort             int          `default:"4000"`
	Clients             []Client     `required:"true"`
	Secret              string       `required:"true"`
	IsPrefork           bool         `default:"false"`
	IsProd              bool         `default:"false"`
	BrokerTCPUrl        string       `requred:"true"`
	Subscribers         []Subscriber `required:"false"`
	StoragePath         string       `default:"./data"`
	InfluxDbAccessToken string       `requred:"true"`
	InfluxDbURL         string       `requred:"true"`
	InfluxDbOrg         string       `requred:"true"`
	InfluxDbBucket      string       `requred:"true"`
	StorageWorkersCount int          `default:"2"`
	PartitionDuration   int64        `default:"1"`
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
