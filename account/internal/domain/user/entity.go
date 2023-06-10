package user

import (
	"github.com/shopspring/decimal"
	"time"
)

type Entity struct {
	ID           string           `json:"id" db:"id"`
	UserName     *string          `json:"user_name" db:"user_name"`
	Type         *int             `json:"type" db:"account_type_id"`
	Latitude     *decimal.Decimal `json:"latitude,omitempty" db:"latitude" validate:"required"`
	Longitude    *decimal.Decimal `json:"longitude,omitempty" db:"latitude" validate:"required"`
	PhoneNumber  *string          `json:"phone_number" db:"phone_number"`
	ProfilePhoto *string          `json:"profile_photo" db:"profile_photo"`
	CreatedAt    time.Time        `json:"-" db:"created_at"`
	UpdatedAt    time.Time        `json:"-" db:"updated_at"`
}

type Types struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}
