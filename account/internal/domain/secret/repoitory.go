package secret

import (
	"context"
	"github.com/KbtuGophers/kiosk/account/internal/domain/activity"
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, data Entity) error
	GetByKey(ctx context.Context, key string) (Entity, error)
	DeleteExpiredTokens(otpInterval int) error
	Lock(ctx context.Context, key string) (func(), *sqlx.Tx, error)
	Update(ctx context.Context, tx *sqlx.Tx, key string, data *UpdateRequest) error
	GetAccountByPhone(phone string) (user.Entity, error)
	CheckForActivities(data activity.Entity) error
}
