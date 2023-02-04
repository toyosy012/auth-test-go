package models

import (
	"time"
)

func NewAccessTokenInput(accountID, email string, now, expiration time.Time) AccessTokenInput {
	return AccessTokenInput{
		accountID: accountID,
		email:     email,
		now:       now,
		expiredAt: expiration,
	}
}

// AccessTokenInput TODO Register Claim NamesとPrivate Claim Namesを別途定義して組み込むべき?
type AccessTokenInput struct {
	accountID string
	email     string
	now       time.Time
	expiredAt time.Time
}

func (i AccessTokenInput) AccountID() string    { return i.accountID }
func (i AccessTokenInput) Email() string        { return i.email }
func (i AccessTokenInput) Now() time.Time       { return i.now }
func (i AccessTokenInput) ExpiredAt() time.Time { return i.expiredAt }

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

func NewToken(access, refresh string) Token { return Token{access: access, refresh: refresh} }

type Token struct {
	access  string
	refresh string
}

func (t Token) Access() string  { return t.access }
func (t Token) Refresh() string { return t.refresh }

func NewTokenOwner(id string, email string) TokenOwner { return TokenOwner{id: id, email: email} }

type TokenOwner struct {
	id    string
	email string
}

func (o TokenOwner) ID() string    { return o.id }
func (o TokenOwner) Email() string { return o.email }
