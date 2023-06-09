package account

import (
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
)

type Configuration func(s *Service) error

type Service struct {
	accountRepository user.Repository

	accountCache user.Cache
}

func New(configs ...Configuration) (s *Service, err error) {
	// Create the service
	s = &Service{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the service into the configuration function
		if err = cfg(s); err != nil {
			return
		}
	}
	return
}

func WithAccountRepository(accountRepository user.Repository) Configuration {
	return func(s *Service) error {
		s.accountRepository = accountRepository
		return nil
	}
}

func WithAccountCache(accountCache user.Cache) Configuration {
	return func(s *Service) error {
		s.accountCache = accountCache
		return nil
	}
}
