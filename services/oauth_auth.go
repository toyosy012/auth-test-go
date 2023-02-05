package services

import (
	"time"

	"auth-test/models"
)

type OauthToken interface {
	Claim(string, string, string, time.Time) (*models.Token, error)
	Refresh(string, string, time.Time) (*models.Token, error)
	Verify(string) error
}

func NewTokenAuthorization(
	authorizer models.Authorizer,
	tokenRepo models.TokenAccessor,
	userAccountRepo models.UserAccountAccessor,
	refreshExpiration time.Duration,
	accessExpiration time.Duration,
) TokenAuthorization {
	return TokenAuthorization{
		authorizer:        authorizer,
		tokenRepo:         tokenRepo,
		userAccountRepo:   userAccountRepo,
		refreshExpiration: refreshExpiration,
		accessExpiration:  accessExpiration,
	}
}

type TokenAuthorization struct {
	authorizer        models.Authorizer
	tokenRepo         models.TokenAccessor
	userAccountRepo   models.UserAccountAccessor
	refreshExpiration time.Duration
	accessExpiration  time.Duration
}

func (a TokenAuthorization) Claim(email, password, newRefreshToken string, now time.Time) (*models.Token, error) {
	account, err := a.userAccountRepo.FindByEmail(email)
	if err != nil {
		return nil, NewApplicationErr(FailedCreateToken, err)
	}

	hash := models.NewEncryptedPassword(account.Password())
	if err = hash.MatchWith(password); err != nil {
		return nil, NewApplicationErr(FailedCreateToken, err)
	}

	refreshToken, err := a.tokenRepo.Insert(models.NewRefreshTokenInput(
		account.ID(), newRefreshToken, now.Add(a.refreshExpiration),
	))
	if err != nil {
		return nil, NewApplicationErr(FailedCreateToken, err)
	}

	accessToken, err := a.authorizer.Sign(models.NewAccessTokenInput(
		account.ID(), account.Email(), now, now.Add(a.accessExpiration),
	))
	if err != nil {
		return nil, NewApplicationErr(FailedCreateToken, err)
	}

	response := models.NewToken(accessToken, refreshToken)
	return &response, nil
}

func (a TokenAuthorization) Refresh(newRefreshToken, oldRefreshToken string, now time.Time) (*models.Token, error) {
	owner, err := a.tokenRepo.FindOwner(oldRefreshToken, now.UTC())
	if err != nil {
		return nil, NewApplicationErr(FailedCreateToken, err)
	}

	refreshToken, err := a.tokenRepo.Insert(models.NewRefreshTokenInput(
		owner.ID(), newRefreshToken, now.Add(a.refreshExpiration),
	))
	if err != nil {
		return nil, NewApplicationErr(FailedCreateToken, err)
	}

	accessToken, err := a.authorizer.Sign(models.NewAccessTokenInput(
		owner.ID(), owner.Email(), now, now.Add(a.accessExpiration),
	))
	if err != nil {
		return nil, NewApplicationErr(FailedCreateToken, err)
	}

	response := models.NewToken(accessToken, refreshToken)
	return &response, nil
}

func (a TokenAuthorization) Verify(accessToken string) error {
	if err := a.authorizer.Verify(accessToken); err != nil {
		return NewApplicationErr(FailedAuthenticate, err)
	}

	return nil
}
