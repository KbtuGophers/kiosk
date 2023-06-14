package otp

import (
	"context"
	"errors"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/activity"
	"github.com/KbtuGophers/kiosk/account/internal/domain/secret"
	"github.com/KbtuGophers/kiosk/account/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/xlzd/gotp"
	"time"
)

/*

	Create(data *Request) error
	GetByKey(key string) (*Response, error)
	DeleteExpiredTokens(otpInterval string)
	Lock(ctx context.Context, key string) (func(), *sqlx.Tx, error)
	Update(ctx context.Context, tx *sqlx.Tx, key string, data *UpdateRequest) error


*/

func (s *Service) Create(ctx context.Context, req secret.Request) (res secret.Response, err error) {

	_, err = s.GetAccountByPhone(req.PhoneNumber)
	if err != nil {
		return
	}

	OtpKey := uuid.New().String()
	OtpSecret := gotp.RandomSecret(16)
	code := gotp.NewTOTP(OtpSecret, 4, s.OtpInterval, nil).Now()

	data := secret.Entity{
		ID:          uuid.New().String(),
		Key:         &OtpKey,
		Secret:      &OtpSecret,
		Status:      1,
		PhoneNumber: &req.PhoneNumber,
		SendAt:      time.Now().Unix(),
		DebugMode:   req.DebugMode,
	}

	err = s.OtpRepository.Create(ctx, data)
	if err != nil {
		return
	}

	if req.DebugMode == true {
		data.Code = &code

	} else {
		//params := &openapi.CreateMessageParams{}
		//params.SetTo(req.PhoneNumber)
		//params.SetBody("Your code: " + code)
		//
		//_, err = s.client.Api.CreateMessage(params)
		//if err != nil {
		//	return
		//}

		message := fmt.Sprintf("Пароль авторизации:%s", code)
		//fmt.Println(message)
		_, err = s.client.SendSms(message, req.PhoneNumber)
		if err != nil {
			return
		}
	}
	res = secret.ParseFromEntity(data)

	return

}

func (s *Service) GetByKey(ctx context.Context, key string) (res secret.Response, err error) {
	//fmt.Println("GetByKey: ", s.OtpRepository)
	data, err := s.OtpRepository.GetByKey(ctx, key)
	if err != nil {
		//fmt.Println("errorrrrr")
		return
	}

	res = secret.ParseFromEntity(data)
	return
}

func (s *Service) Check(ctx context.Context, req secret.Request) (res secret.Response, err error) {
	reqTime := time.Now().Unix()
	updReq := &secret.UpdateRequest{}
	//fmt.Println("Key: ", req.Key)
	var data secret.Entity
	data, err = s.OtpRepository.GetByKey(ctx, req.Key)
	if err != nil {
		return
	}
	//check status
	if data.Status != 1 {
		err = errors.New("otp is invalid")
		return
	}

	//check if time expired
	if reqTime-data.SendAt > int64(s.OtpInterval) {
		updReq.Status = 0
		err = errors.New("otp time is expired")
		return
	}

	//check attempts
	data.Attempts += 1
	if data.Attempts > s.OtpAttempts {
		updReq.Status = 0
		err = errors.New(fmt.Sprintf("attempted %s times", data.Attempts))
		return
	}

	valid := gotp.NewTOTP(*data.Secret, 4, s.OtpInterval, nil).Verify(req.Code, data.SendAt)
	if !valid {
		err = errors.New("code is invalid")
		return
	} else {
		updReq.Status = 0
		updReq.ConfirmedAt = time.Now().Unix()
	}

	if data.Attempts < s.OtpAttempts {
		updReq.Attempts = data.Attempts

		var unlock func()
		var reqTX *sqlx.Tx
		unlock, reqTX, err = s.OtpRepository.Lock(ctx, *data.Key)
		if err != nil {
			return
		}
		defer unlock()
		err = s.OtpRepository.Update(ctx, reqTX, *data.Key, updReq)

		if err != nil {
			return
		}

	}

	res = secret.ParseFromEntity(data)

	//fmt.Println("GetById is finished")
	return

}

func (s *Service) InsertActivities(accountId string) error {
	data := activity.Entity{Activity: "LOGIN", AccountId: accountId}
	if err := s.OtpRepository.CheckForActivities(data); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAccountByPhone(phone string) (res user.Response, err error) {
	var data user.Entity
	data, err = s.OtpRepository.GetAccountByPhone(phone)
	if err != nil {
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) DeleteExpiredTokens() (err error) {
	err = s.OtpRepository.DeleteExpiredTokens(s.OtpInterval)
	if err != nil {
		fmt.Println("___________erire_______________" + err.Error())
	}
	return
}
