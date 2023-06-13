package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/activity"
	"github.com/KbtuGophers/kiosk/account/internal/domain/secret"
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type OtpRepository struct {
	db *sqlx.DB
}

func (o *OtpRepository) Create(ctx context.Context, data secret.Entity) error {

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

func (o *OtpRepository) DeleteExpiredTokens(otpInterval int) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = o.db.NamedExecContext(ctx, `
		DELETE FROM otps
		WHERE id IN (
		    SELECT id
		    FROM otps
		    WHERE created_at < (CURRENT_TIMESTAMP - INTERVAL '`+strconv.Itoa(otpInterval)+` seconds')
		)`,
		map[string]interface{}{})

	return
}

func (o *OtpRepository) Lock(ctx context.Context, key string) (func(), *sqlx.Tx, error) {
	query := `SELECT * FROM otps WHERE key=:key FOR UPDATE`
	args := map[string]interface{}{
		"key": key,
	}

	tx, err := o.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	return func() {
		err = tx.Commit()
		if err != nil {
			return
		}

	}, tx, nil

}

func (o *OtpRepository) Update(ctx context.Context, tx *sqlx.Tx, key string, data *secret.UpdateRequest) error {
	var args []any
	var setValues []string
	ind := 1
	if data.Status != 1 {
		setValues = append(setValues, fmt.Sprintf("status=$%d", ind))
		args = append(args, data.Status)
		ind++
	}

	if data.Attempts != 0 {
		setValues = append(setValues, fmt.Sprintf("attempts=$%d", ind))
		args = append(args, data.Attempts)
		ind++
	}

	if data.ConfirmedAt != 0 {
		setValues = append(setValues, fmt.Sprintf("confirmed_at=$%d", ind))
		args = append(args, data.ConfirmedAt)
		ind++
	}

	if data.SendAt != 0 {
		setValues = append(setValues, fmt.Sprintf("send_at=$%d", ind))
		args = append(args, data.SendAt)
		ind++
	}

	//args = append(args, key)
	setValue := strings.Join(setValues, ",")

	query := fmt.Sprintf("UPDATE otps SET %s WHERE key='%s'", setValue, key)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	fmt.Println("success")
	return nil

}

func (o *OtpRepository) GetAccountByPhone(phone string) (user.Entity, error) {
	var data user.Entity

	query := `
		SELECT *
		FROM accounts
		WHERE phone_number=$1 LIMIT 1`

	args := []any{phone}

	err := o.db.Get(&data, query, args...)
	if err != nil {

		return user.Entity{}, err
	}

	return data, nil
}

func (o *OtpRepository) CheckForActivities(data activity.Entity) error {
	query := `
		INSERT INTO user_activities (account_id, activity, timestamp) VALUES ($1, $2, $3)
    `
	args := []any{data.AccountId, data.Activity, data.Timestamp}
	if _, err := o.db.Exec(query, args...); err != nil {
		return err
	}

	return nil

}

func NewOtpRepository(db *sqlx.DB) *OtpRepository {
	return &OtpRepository{db: db}
}
