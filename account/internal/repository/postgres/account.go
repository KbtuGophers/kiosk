package postgres

import (
	"context"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func (a *AccountRepository) Select(ctx context.Context) (dest []user.Entity, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepository) Create(ctx context.Context, data user.Entity) (id string, err error) {
	fmt.Println(data)
	query := `
		INSERT INTO account.accounts (id, user_name, phone_number, account_type_id, profile_photo) 
		VALUES ($1, $2, $3, $4, $5)
	`

	args := []any{data.ID, data.UserName, data.PhoneNumber, data.Type, data.ProfilePhoto}

	if _, err = a.db.ExecContext(ctx, query, args...); err != nil {
		return
	}

	id = data.ID

	return
}

func (a *AccountRepository) Get(ctx context.Context, id string) (dest user.Entity, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepository) Update(ctx context.Context, id string, data user.Entity) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepository) Delete(ctx context.Context, id string) (err error) {
	//TODO implement me
	panic("implement me")
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}
