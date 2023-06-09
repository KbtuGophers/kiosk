package user

import (
	"github.com/shopspring/decimal"
	"time"
)

type Entity struct {
	ID           string           `json:"id"`
	UserName     *string          `json:"user_name"`
	Type         *int             `json:"type"`
	Latitude     *decimal.Decimal `json:"latitude,omitempty" validate:"required"`
	Longitude    *decimal.Decimal `json:"longitude,omitempty" validate:"required"`
	PhoneNumber  *string          `json:"phone_number"`
	ProfilePhoto *string          `json:"profile_photo"`
	CreatedAt    time.Time        `json:"-"`
	UpdatedAt    time.Time        `json:"-"`
}

type Types struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}
