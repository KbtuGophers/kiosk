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

func (s *Service) GetAuthor(ctx context.Context, id string) (res user.Response, err error) {
	data, err := s.accountRepository.Get(ctx, id)
	if err != nil {
		return
	}
	res = user.ParseFromEntity(data)

	return
}

func (s *Service) UpdateAuthor(ctx context.Context, id string, req user.Request) error {
	data := user.Entity{
		UserName:    &req.UserName,
		Type:        &req.Type,
		Longitude:   &req.Longitude,
		Latitude:    &req.Latitude,
		PhoneNumber: &req.PhoneNumber,
	}
	return s.accountRepository.Update(ctx, id, data)
}

func (s *Service) DeleteAuthor(ctx context.Context, id string) error {
	return s.accountRepository.Delete(ctx, id)
}
