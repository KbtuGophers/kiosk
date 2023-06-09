package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/secret"
	"github.com/jmoiron/sqlx"
)

type OtpRepository struct {
	db *sqlx.DB
}

func (o *OtpRepository) Create(ctx context.Context, data secret.Entity) error {
	fmt.Println("data to create: ", o.db)
	query := `
	INSERT INTO otps 
    (id, key, secret, phone_number, send_at, attempts, status, confirmed_at) 
	VALUES ($1,$2, $3, $4, $5, $6, $7, $8)
	`

	args := []any{data.ID, data.Key, data.Secret, data.PhoneNumber, data.SendAt, data.Attempts, data.Status, data.ConfirmedAt}

	if _, err := o.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil

}

func (o *OtpRepository) GetByKey(ctx context.Context, key string) (data secret.Entity, err error) {
	//var rows []secret.Entity

	query := `
		SELECT *
		FROM account.otps
		WHERE key=$1 LIMIT 1`

	args := []any{key}

	if err = o.db.GetContext(ctx, &data, query, args...); err != nil && err != sql.ErrNoRows {
		return
	}
	//data = rows[0]

	err = nil

	return
}

func (o *OtpRepository) DeleteExpiredTokens(otpInterval string) {
	query := `
		DELETE FROM opts WHERE created_at < (CURRENT_TIMESTAMP - INTERVAL ` + otpInterval + ` seconds)
    `
	fmt.Println(query)
	return
}

func (o *OtpRepository) Lock(ctx context.Context, key string) (func(), *sqlx.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OtpRepository) Update(ctx context.Context, tx *sqlx.Tx, key string, data secret.UpdateRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewOtpRepository(db *sqlx.DB) *OtpRepository {
	return &OtpRepository{db: db}
}
