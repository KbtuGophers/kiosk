package account

import (
	"context"
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/google/uuid"
)

func (s *Service) AddAccount(ctx context.Context, req user.Request) (res user.Response, err error) {

	data := user.Entity{
		ID:           uuid.New().String(),
		UserName:     &req.UserName,
		Type:         &req.Type,
		Latitude:     &req.Latitude,
		Longitude:    &req.Longitude,
		PhoneNumber:  &req.PhoneNumber,
		ProfilePhoto: &req.ProfilePhoto,
	}
	data.ID, err = s.accountRepository.Create(ctx, data)
	if err != nil {
		return
	}

	res = user.ParseFromEntity(data)
	return
}

func (s *Service) GetAccount(ctx context.Context, id string) (res user.Response, err error) {
	data, err := s.accountRepository.Get(ctx, id)
	if err != nil {
		return
	}
	res = user.ParseFromEntity(data)

	return
}

func (s *Service) UpdateAccount(ctx context.Context, id string, req user.Request) error {
	data := user.ParseFromRequest(req)
	//fmt.Println("-----------------------------------")
	//fmt.Println(req.Longitude.String())
	//fmt.Println("-----------------------------------")

	if err := s.accountRepository.Update(ctx, id, data); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteAccount(ctx context.Context, id string) error {
	return s.accountRepository.Delete(ctx, id)
}

func (s *Service) GetAllAccounts(ctx context.Context) (res []user.Response, err error) {
	data, err := s.accountRepository.List(ctx)
	if err != nil {
		return
	}

	res = user.ParseFromEntities(data)

	return
}

func (s *Service) CreateAccountType(accountType user.Types) (int, error) {
	return s.accountRepository.InsertType(accountType)
}
