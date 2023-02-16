package dto

import "time"

type SubscriberInput struct {
	ApiaryId string `json:"apiaryId" required:"true"`
	HiveId   string `json:"hiveId" required:"true"`
}

type SubscriberOutput struct {
	ApiaryId  string    `json:"apiaryId"`
	HiveId    string    `json:"hiveId"`
	CreatedAt time.Time `json:"createdAt"`
}
