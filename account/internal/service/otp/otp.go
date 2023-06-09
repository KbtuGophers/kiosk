package otp

import (
	"context"
	"errors"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/secret"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/xlzd/gotp"
	"strconv"
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
	}

	err = s.OtpRepository.Create(ctx, data)
	if err != nil {
		return
	}

	if req.DebugMode == true {
		data.Code = &code
	} else {
		params := &openapi.CreateMessageParams{}
		params.SetTo(req.PhoneNumber)
		params.SetBody("Your code: " + code)

		_, err = s.client.Api.CreateMessage(params)
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
		fmt.Println("errorrrrr")
		return
	}

	res = secret.ParseFromEntity(data)
	return
}

func (s *Service) GetById(ctx context.Context, req secret.Request) (res secret.Response, err error) {
	reqTime := time.Now().Unix()
	updReq := secret.UpdateRequest{}

	//fmt.Println("Key: ", req.Key)

	res, err = s.GetByKey(ctx, req.Key)
	if err != nil {
		return
	}
	//check status
	if res.Status != 1 {
		err = errors.New("otp is invalid")
		return
	}

	//check if time expired
	if reqTime-res.SendAt > int64(s.OtpInterval) {
		updReq.Status = 0
		err = errors.New("otp time is expired")
		return
	}

	//check attempts
	res.Attempts += 1
	if res.Attempts > s.OtpAttempts {
		updReq.Status = 0
		err = errors.New(fmt.Sprintf("attempted %s times", res.Attempts))
		return
	}

	valid := gotp.NewTOTP(res.Secret, 4, s.OtpInterval, nil).Verify(req.Code, res.SendAt)
	if !valid {
		err = errors.New("code is invalid")
		return
	} else {
		updReq.Status = 0
		updReq.ConfirmedAt = time.Now().Unix()
	}

	if res.Attempts < s.OtpAttempts {
		updReq.Attempts = res.Attempts

		var unlock func()
		var reqTX *sqlx.Tx
		unlock, reqTX, err = s.OtpRepository.Lock(ctx, res.Key)
		if err != nil {
			return
		}
		defer unlock()
		err = s.OtpRepository.Update(ctx, reqTX, res.Key, updReq)
		if err != nil {
			return
		}
	}

	fmt.Println("GetById is finished")
	return

}

func (s *Service) DeleteExpiredTokens(ctx context.Context) {
	s.OtpRepository.DeleteExpiredTokens(strconv.Itoa(s.OtpInterval))
}
