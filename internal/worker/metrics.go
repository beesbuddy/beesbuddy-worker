package worker

type metrics struct {
	ClientId    string `json:"clientId"`
	ApiaryId    string `json:"apiaryId"`
	HiveId      string `json:"hiveId"`
	Temperature string `json:"t"`
	Humidity    string `json:"h"`
	Weight      string `json:"w"`
}
