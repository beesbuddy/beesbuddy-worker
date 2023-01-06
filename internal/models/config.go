package models

type Config struct {
	AppName        string   `default:"BeesBuddy Worker"`
	AppHost        string   `default:"0.0.0.0"`
	AppPort        int      `default:"4000"`
	Clients        []Client `required:"true"`
	Secret         string   `required:"true"`
	UiHotReloadUrl string   `default:"http://localhost:8080"`
	IsPrefork      bool     `default:"false"`
	IsProd         bool     `default:"false"`
	BrokerTCPUrl   string   `requred:"true"`
}

type Client struct {
	SecretKey string `required:"true"`
	AppKey    string `required:"true"`
}
