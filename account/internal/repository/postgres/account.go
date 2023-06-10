package postgres

import (
	"context"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/jmoiron/sqlx"
	"strings"
)

type AccountRepository struct {
	db *sqlx.DB
}

func (a *AccountRepository) Select(ctx context.Context) ([]user.Entity, error) {
	var dest []user.Entity
	err := a.db.SelectContext(ctx, &dest, "SELECT * FROM accounts", nil)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (a *AccountRepository) Create(ctx context.Context, data user.Entity) (id string, err error) {
	fmt.Println(data)
	query := `
		INSERT INTO accounts (id, user_name, phone_number, account_type_id, profile_photo) 
		VALUES ($1, $2, $3, $4, $5)
	`

	args := []any{data.ID, data.UserName, data.PhoneNumber, data.Type, data.ProfilePhoto}

	if _, err = a.db.ExecContext(ctx, query, args...); err != nil {
		return
	}

	id = data.ID

	return
}

func (a *AccountRepository) Get(ctx context.Context, id string) (user.Entity, error) {
	query := "SELECT * FROM accounts WHERE id=$1"
	var dest user.Entity
	args := []any{id}

	if err := a.db.GetContext(ctx, &dest, query, args...); err != nil {
		return user.Entity{}, err
	}

	return dest, nil

}

func (a *AccountRepository) Update(ctx context.Context, id string, data user.Entity) (err error) {
	sets, args := a.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE accounts SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		_, err = a.db.ExecContext(ctx, query, args...)
		if err != nil {
			return
		}

	}

	return
}

func (a *AccountRepository) prepareArgs(data user.Entity) (sets []string, args []any) {
	if data.UserName != nil {
		args = append(args, data.UserName)
		sets = append(sets, fmt.Sprintf("user_name=$%d", len(args)))
	}

	//if data.Latitude != nil {
	//	args = append(args, data.Longitude)
	//	sets = append(sets, fmt.Sprintf("longtitude=$%d", len(args)))
	//}
	//
	//if data.Latitude != nil {
	//	args = append(args, data.Latitude)
	//	sets = append(sets, fmt.Sprintf("latitude=$%d", len(args)))
	//}
	if data.PhoneNumber != nil {
		args = append(args, data.PhoneNumber)
		sets = append(sets, fmt.Sprintf("phone_number=$%d", len(args)))
	}
	if data.ProfilePhoto != nil {
		args = append(args, data.ProfilePhoto)
		sets = append(sets, fmt.Sprintf("profile_photo=$%d", len(args)))
	}

	return
}

func (a *AccountRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE 
		FROM accounts
		WHERE id=$1`

	args := []any{id}

	_, err = a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return
	}

	return
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}
