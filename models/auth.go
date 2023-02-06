package models

import (
	"time"
)

type Authorizer interface {
	Sign(IDTokenInput) (string, error)
	Verify(string) error
}

func NewAccessTokenInput(accountID, email string, now, expiration time.Time) IDTokenInput {
	return IDTokenInput{
		accountID: accountID,
		email:     email,
		now:       now,
		expiredAt: expiration,
	}
}

// IDTokenInput TODO Register Claim NamesとPrivate Claim Namesを別途定義して組み込むべき?
type IDTokenInput struct {
	accountID string
	email     string
	now       time.Time
	expiredAt time.Time
}

func (i IDTokenInput) AccountID() string    { return i.accountID }
func (i IDTokenInput) Email() string        { return i.email }
func (i IDTokenInput) Now() time.Time       { return i.now }
func (i IDTokenInput) ExpiredAt() time.Time { return i.expiredAt }

type TokenAccessor interface {
	Insert(RefreshTokenInput) (string, error)
	Delete(string) error
	FindOwner(string, time.Time) (*TokenOwner, error)
}

func NewRefreshTokenInput(accountID, value string, expiration time.Time) RefreshTokenInput {
	return RefreshTokenInput{
		accountID: accountID,
		value:     value,
		expiredAt: expiration,
	}
}

type RefreshTokenInput struct {
	accountID string
	expiredAt time.Time
	value     string
}

func (i RefreshTokenInput) AccountID() string    { return i.accountID }
func (i RefreshTokenInput) Value() string        { return i.value }
func (i RefreshTokenInput) ExpiredAt() time.Time { return i.expiredAt }

func NewToken(id, refresh string) Token { return Token{idToken: id, refresh: refresh} }

type Token struct {
	idToken string
	refresh string
}

func (t Token) IDToken() string { return t.idToken }
func (t Token) Refresh() string { return t.refresh }

func NewTokenOwner(id string, email string) TokenOwner { return TokenOwner{id: id, email: email} }

type TokenOwner struct {
	id    string
	email string
}

func (o TokenOwner) ID() string    { return o.id }
func (o TokenOwner) Email() string { return o.email }
