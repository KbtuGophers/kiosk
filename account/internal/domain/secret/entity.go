package secret

import "time"

type Entity struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"-" db:"phone_number"`
	UpdatedAt   time.Time `json:"-" db:"phone_number"`
	PhoneNumber *string   `json:"phone_number" db:"phone_number"`
	Key         *string   `json:"key"`
	Secret      *string   `json:"secret"`
	Code        *string   `json:"debug_code"`
	Attempts    int       `json:"-"`
	Status      int       `json:"-"`
	SendAt      int64     `json:"-"`
	ConfirmedAt int64     `json:"-"`
	DebugMode   bool      `json:"debug_mode"`
}
