package repository

import (
	"github.com/KbtuGophers/kiosk/account/internal/domain/secret"
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/KbtuGophers/kiosk/account/internal/repository/postgres"
	"github.com/KbtuGophers/kiosk/account/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres *store.Database
	Account  user.Repository
	Otp      secret.Repository
}

const (
	accountTable     = "account.accounts"
	accountTypeTable = "account.account_types"
	otpTable         = "account.otps"
)

func New(configs ...Configuration) (r *Repository, err error) {
	// Create the repository
	r = &Repository{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the repository into the configuration function
		if err = cfg(r); err != nil {
			return
		}
	}

	return
}

func (r *Repository) Close() {
	if r.postgres != nil {
		r.postgres.Client.Close()
	}
}

func WithPostgresStore(schema, dataSourceName string) Configuration {
	return func(r *Repository) (err error) {
		r.postgres, err = store.NewDatabase(schema, dataSourceName)
		if err != nil {
			return
		}

		if err = r.postgres.Migrate(); err != nil && err.Error() != "no change" {
			return
		}
		err = nil

		r.Account = postgres.NewAccountRepository(r.postgres.Client)
		r.Otp = postgres.NewOtpRepository(r.postgres.Client)

		return
	}
}
