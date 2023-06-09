package service

import (
	"github.com/KbtuGophers/kiosk/account/internal/service/account"
	"github.com/KbtuGophers/kiosk/account/internal/service/otp"
)

type Service struct {
	AccountService account.Service
	OtpService     otp.Service
}

func NewService(AccountService account.Service, OtpService otp.Service) *Service {
	return &Service{
		AccountService: AccountService,
		OtpService:     OtpService,
	}
}
