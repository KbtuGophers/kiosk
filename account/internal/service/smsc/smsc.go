package smsc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KbtuGophers/kiosk/account/internal/domain/sms"
	"net/http"
	"net/url"
	"strconv"
)

func (s *Client) Send(params *url.Values) (*sms.Response, error) {

	//fmt.Println(params.Encode())

	req, err := http.NewRequest("POST", "https://smsc.kz/sys/send.php?", bytes.NewBufferString(params.Encode()))

	//fmt.Println(bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &sms.Response{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	return res, nil
}

func (s *Client) SendSms(message string, phones ...string) (*sms.Response, error) {
	s.mu.RLock()

	if len(phones) == 0 {
		s.mu.RUnlock()
		return nil, errors.New("did not set phone numbers")
	}

	var phones_str = ""
	var sep = ""
	for _, phone := range phones {
		phones_str += sep + phone
		sep = ","
	}

	params := url.Values{}
	params.Add("login", s.login)
	params.Add("psw", s.password)
	params.Add("phones", phones_str)
	params.Add("mes", message)
	params.Add("charset", s.charset)
	params.Add("fmt", "3") // json
	params.Add("tinyurl", "1")
	if s.sender != "" {
		params.Add("sender", s.sender)
	}

	s.mu.RUnlock()

	fmt.Println(params.Encode())

	return s.Send(&params)
}
