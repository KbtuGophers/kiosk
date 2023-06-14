package activity

import "time"

type Entity struct {
	ID        int    `json:"id"`
	AccountId string `json:"account_id"`
	Activity  string `json:"activity"`
	Timestamp time.Time
}
