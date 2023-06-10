package user

import (
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type Request struct {
	UserName     string          `json:"user_name"`
	Type         int             `json:"type"`
	Latitude     decimal.Decimal `json:"latitude,omitempty" validate:"required"`
	Longitude    decimal.Decimal `json:"longitude,omitempty" validate:"required"`
	PhoneNumber  string          `json:"phone_number"`
	ProfilePhoto string          `json:"profile_photo"`
	CreatedAt    time.Time       `json:"-"`
	UpdatedAt    time.Time       `json:"-"`
}

func (s *Request) Bind(r *http.Request) error {
	//if s.UserName == "" {
	//	return errors.New("username: cannot be blank")
	//}

	//if s.Type == 0 {
	//	return errors.New("account type: cannot be blank")
	//}

	//if s.PhoneNumber == "" {
	//	return errors.New("phone number: cannot be blank")
	//}

	return nil
}

type Response struct {
	ID           string          `json:"id"`
	UserName     string          `json:"user_name"`
	Type         int             `json:"type"`
	Latitude     decimal.Decimal `json:"latitude,omitempty"`
	Longitude    decimal.Decimal `json:"longitude,omitempty"`
	PhoneNumber  string          `json:"phone_number"`
	ProfilePhoto string          `json:"profile_photo"`
	CreatedAt    time.Time       `json:"-" db:"created_at"`
	UpdatedAt    time.Time       `json:"-"`
}

func ParseFromEntity(data Entity) (res Response) {

	res = Response{
		ID:           data.ID,
		UserName:     *data.UserName,
		Type:         *data.Type,
		PhoneNumber:  *data.PhoneNumber,
		ProfilePhoto: *data.ProfilePhoto,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}

	if data.Latitude != nil {
		res.Latitude = *data.Latitude
	}
	if data.Longitude != nil {
		res.Longitude = *data.Longitude
	}

	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
