package settings

import "time"

type AppSettings struct {
	AppName      string       `default:"BeesBuddy Worker"`
	AppHost      string       `default:"0.0.0.0"`
	AppPort      int          `default:"4000"`
	Clients      []Client     `required:"true"`
	Secret       string       `required:"true"`
	IsPrefork    bool         `default:"false"`
	IsProd       bool         `default:"false"`
	BrokerTCPUrl string       `requred:"true"`
	Subscribers  []Subscriber `required:"false"`
	DbName       string       `default:"beesbuddy"`
}

type Client struct {
	AppKey string `required:"true"`
}

type Subscriber struct {
	ApiaryId  string    `json:"apiaryId" required:"true"`
	HiveId    string    `json:"hiveId" required:"true"`
	CreatedAt time.Time `json:"createdAt" required:"true"`
}