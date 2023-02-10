package dto

type SubscriberInput struct {
	ApiaryId string `json:"apiaryId" required:"true"`
	HiveId   string `json:"hiveId" required:"true"`
}
