package smsc

import (
	"errors"
	"net/http"
	"sync"
)

type Client struct {
	client *http.Client

	mu       sync.RWMutex
	login    string
	password string
	tinyurl  string
	charset  string
	sender   string
}

func New(login, password, sender string) (*Client, error) {
	if login == "" || password == "" {
		return nil, errors.New("empty login or password")
	}

	sc := &Client{
		login:    login,
		password: password,
		charset:  "utf-8",
		client:   &http.Client{},
		sender:   sender,
	}

	return sc, nil
}
