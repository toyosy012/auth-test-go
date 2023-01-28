package services

import (
	"time"

	"github.com/google/uuid"

	"auth-test/models"
)

type Login interface {
	Sign(string, string) (string, error)
	Verify(string) error
}

func NewAuthorizer(userAccountRepo models.UserAccountRepository, authorizer models.Authorizer) Authorizer {
	return Authorizer{
		userAccountRepo: userAccountRepo,
		authorizer:      authorizer,
	}
}

type Authorizer struct {
	userAccountRepo models.UserAccountRepository
	authorizer      models.Authorizer
}

func (a Authorizer) Sign(email, password string) (string, error) {
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

func (a Authorizer) Verify(token string) error {
	return a.authorizer.Verify(token)
}

type StoredAuthorizer interface {
	Sign(string, string) (string, error)
	Verify(string) error
	FindUser(string, string) error
	SignOut(string) error
}

func NewStoredAuthorization(
	a models.UserAccountRepository,
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
	userAccountRepo  models.UserAccountRepository
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

func (a StoredAuthorization) SignOut(token string) error {
	return a.userSessionRepo.Delete(token)
}
