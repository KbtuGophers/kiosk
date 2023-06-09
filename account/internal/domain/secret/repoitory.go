package secret

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, data Entity) error
	GetByKey(ctx context.Context, key string) (Entity, error)
	DeleteExpiredTokens(otpInterval string)
	Lock(ctx context.Context, key string) (func(), *sqlx.Tx, error)
	Update(ctx context.Context, tx *sqlx.Tx, key string, data UpdateRequest) error
}
