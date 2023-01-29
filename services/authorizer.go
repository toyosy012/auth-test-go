package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"auth-test/models"
)

type Authorizer interface {
	Sign(string, string) (string, error)
	Verify(string) error
}

func NewAuthorizer(userAccountRepo models.UserAccountAccessor, authorizer models.Authorizer) Authorization {
	return Authorization{
		userAccountRepo: userAccountRepo,
		authorizer:      authorizer,
	}
}

type Authorization struct {
	userAccountRepo models.UserAccountAccessor
	authorizer      models.Authorizer
}

func (a Authorization) Sign(email, password string) (string, error) {
	account, err := a.userAccountRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	encrypted := models.EncryptedPassword{Hash: account.Password}
	if err = encrypted.MatchWith(password); err != nil {
		return "", err
	}

	token, err := a.authorizer.Sign(*account)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Authorization) Verify(token string) error {
	return a.authorizer.Verify(token)
}

type StoredAuthorizer interface {
	Sign(string, string) (string, error)
	Verify(string) error
	FindUser(string, string) error
	SignOut(string, string) error
}

func NewStoredAuthorization(
	a models.UserAccountAccessor,
	s models.UserSessionAccessor,
	availabilityTime time.Duration,
) StoredAuthorization {
	return StoredAuthorization{
		userAccountRepo:  a,
		userSessionRepo:  s,
		availabilityTime: availabilityTime,
	}
}

type StoredAuthorization struct {
	userAccountRepo  models.UserAccountAccessor
	userSessionRepo  models.UserSessionAccessor
	availabilityTime time.Duration
}

func (a StoredAuthorization) Sign(email, password string) (string, error) {
	account, err := a.userAccountRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	encrypted := models.EncryptedPassword{Hash: account.Password}
	if err = encrypted.MatchWith(password); err != nil {
		return "", err
	}

	s := models.NewSession(account.ID, uuid.New().String(), time.Now().UTC().Add(a.availabilityTime))
	return a.userSessionRepo.Register(s)
}

func (a StoredAuthorization) Verify(token string) error {
	err := a.userSessionRepo.Verify(token)
	if err != nil {
		return err
	}

	return nil
}

func (a StoredAuthorization) FindUser(id, token string) error {
	owner, err := a.userSessionRepo.FindUser(token)
	if err != nil {
		return err
	}

	if owner != id {
		return errors.New(http.StatusText(http.StatusNotFound))
	}

	return nil
}

func (a StoredAuthorization) SignOut(owner, token string) error {
	return a.userSessionRepo.Delete(owner, token)
}
