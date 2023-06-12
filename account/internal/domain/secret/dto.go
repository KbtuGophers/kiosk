package secret

import (
	"net/http"
)

type Request struct {
	Key         string `json:"key"`
	Secret      string `json:"secret"`
	Code        string `json:"code"`
	DebugMode   bool   `json:"debug_mode"`
	PhoneNumber string `json:"phone_number"`
}

func (s *Request) Bind(r *http.Request) error {
	//if s.PhoneNumber == "" {
	//	return errors.New("phone number: cannot be blank")
	//}

	return nil
}

type Response struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Code        string `json:"code"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateRequest struct {
	Attempts int   `json:"-"`
	Status   int   `json:"-"`
	SendAt   int64 `json:"-"`

	ConfirmedAt int64 `json:"-"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:  data.ID,
		Key: *data.Key,
	}

	if data.Code != nil {
		res.Code = *data.Code
	}
	if data.Code != nil {
		res.PhoneNumber = *data.PhoneNumber
	}

	return
}
