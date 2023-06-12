package sms

type Response struct {
	Cnt  float64 `json:"cnt"`
	Cost string  `json:"cost"`
	Err  string  `json:"error"`
}
