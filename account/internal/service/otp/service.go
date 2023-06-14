package otp

import (
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/secret"
	"github.com/KbtuGophers/kiosk/account/internal/service/smsc"
)

type Configuration func(otp *Service) error

type Service struct {
	client        *smsc.Client
	OtpAttempts   int
	OtpInterval   int
	OtpRepository secret.Repository
}

func NewOtpService(client *smsc.Client, attempts int, interval int, configs ...Configuration) (s *Service, err error) {
	s = &Service{client: client, OtpAttempts: attempts, OtpInterval: interval}
	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	fmt.Println("OtpService: ", s)

	return
}

func WithOtpRepository(repository secret.Repository) Configuration {
	return func(s *Service) error {
		s.OtpRepository = repository
		return nil
	}
}
