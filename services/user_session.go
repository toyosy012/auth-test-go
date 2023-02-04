package services

import (
	"errors"
	"net/http"
	"time"

	"auth-test/models"
)

type Session interface {
	Sign(string, string, string, time.Time) (string, error)
	Verify(string) error
	FindOwner(string, string) error
	SignOut(string, string) error
}

func NewSessionAuthorization(
	a models.UserAccountAccessor,
	s models.UserSessionAccessor,
	expiration time.Duration,

) UserSession {
	return UserSession{
		userAccountRepo: a,
		userSessionRepo: s,
		expiration:      expiration,
	}
}

type UserSession struct {
	userAccountRepo models.UserAccountAccessor
	userSessionRepo models.UserSessionAccessor
	expiration      time.Duration
}

func (s UserSession) Sign(email, password, sessionID string, now time.Time) (string, error) {
	account, err := s.userAccountRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	hash := models.NewEncryptedPassword(account.Password())
	if err = hash.MatchWith(password); err != nil {
		return "", err
	}

	sess := models.NewSession(account.ID(), sessionID, now.Add(s.expiration))
	return s.userSessionRepo.Register(sess)
}

func (s UserSession) Verify(token string) error {
	err := s.userSessionRepo.Verify(token)
	if err != nil {
		return err
	}

	return nil
}

func (s UserSession) FindOwner(id, token string) error {
	owner, err := s.userSessionRepo.FindOwner(token)
	if err != nil {
		return err
	}

	if owner != id {
		return errors.New(http.StatusText(http.StatusNotFound))
	}

	return nil
}

func (s UserSession) SignOut(owner, token string) error {
	return s.userSessionRepo.Delete(owner, token)
}
