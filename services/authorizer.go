package services

import (
	"time"

	"auth-test/models"
)

func NewAuthorizer(userAccountRepo models.UserAccountRepository, tokenAuthorizer models.TokenAuthorizer) Authorizer {
	return Authorizer{
		userAccountRepo: userAccountRepo,
		tokenAuthorizer: tokenAuthorizer,
	}
}

type Authorizer struct {
	userAccountRepo models.UserAccountRepository
	tokenAuthorizer models.TokenAuthorizer
}

func (a Authorizer) Sign(email, password string, now time.Time) (string, error) {
	account, err := a.userAccountRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	encrypted := models.EncryptedPassword{Hash: account.Password}
	if err = encrypted.MatchWith(password); err != nil {
		return "", err
	}

	token, err := a.tokenAuthorizer.Sign(*account, now)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Authorizer) Verify(token string, now time.Time) error {
	return a.tokenAuthorizer.Verify(token, now)
}
