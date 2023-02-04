package models

import (
	"time"
)

func NewSession(owner, token string, expiredAt time.Time) Session {
	return Session{
		owner:     owner,
		token:     token,
		expiredAt: expiredAt,
	}
}

type Session struct {
	owner     string
	token     string
	expiredAt time.Time
}

func (s Session) Owner() string        { return s.owner }
func (s Session) Token() string        { return s.token }
func (s Session) ExpiredAt() time.Time { return s.expiredAt }

type UserSessionAccessor interface {
	Register(Session) (string, error)
	Verify(string) error
	FindOwner(string) (string, error)
	Delete(string, string) error
}
