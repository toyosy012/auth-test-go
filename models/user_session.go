package models

import (
	"time"
)

type Session struct {
	Owner     string
	Token     string
	ExpiredAt time.Time
}

func NewSession(owner, token string, expiredAt time.Time) Session {
	return Session{
		Owner:     owner,
		Token:     token,
		ExpiredAt: expiredAt,
	}
}

type UserSessionAccessor interface {
	Register(Session) (string, error)
	Verify(string) error
	FindUser(string) (string, error)
	Delete(string, string) error
}
