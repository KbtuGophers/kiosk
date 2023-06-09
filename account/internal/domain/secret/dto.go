package secret

import (
	"net/http"
	"time"
)

type Request struct {
	Key         string `json:"key"`
	Secret      string `json:"secret"`
	Code        string `json:"debug_code"`
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
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	PhoneNumber string    `json:"phone_number"`
	Key         string    `json:"key"`
	Secret      string    `json:"secret"`
	Code        string    `json:"debug_code"`
	Attempts    int       `json:"-"`
	Status      int       `json:"-"`
	SendAt      int64     `json:"-"`
	ConfirmedAt int64     `json:"-"`
	DebugMode   bool      `json:"debug_mode"`
}

type UpdateRequest struct {
	Attempts    int   `json:"-"`
	Status      int   `json:"-"`
	SendAt      int64 `json:"-"`
	ConfirmedAt int64 `json:"-"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:          data.ID,
		Key:         *data.Key,
		Secret:      *data.Secret,
		Code:        *data.Code,
		PhoneNumber: *data.PhoneNumber,
		Status:      data.Status,
		Attempts:    data.Attempts,
		SendAt:      data.SendAt,
		ConfirmedAt: data.ConfirmedAt,
		DebugMode:   data.DebugMode,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return
}
