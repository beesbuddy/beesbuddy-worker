package models

import "time"

type Subscriber struct {
	ApiaryId  string    `json:"apiaryId" required:"true"`
	HiveId    string    `json:"hiveId" required:"true"`
	CreatedAt time.Time `json:"createdAt" required:"true"`
}
