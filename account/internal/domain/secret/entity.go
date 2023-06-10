package secret

import "time"

type Entity struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
	PhoneNumber *string   `json:"phone_number" db:"phone_number"`
	Key         *string   `json:"key"`
	Secret      *string   `json:"secret"`
	Code        *string   `json:"code"`
	Attempts    int       `json:"-"`
	Status      int       `json:"-"`
	SendAt      int64     `json:"-" db:"send_at"`
	ConfirmedAt int64     `json:"-" db:"confirmed_at"`
	DebugMode   bool      `json:"debug_mode"`
}
