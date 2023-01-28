package models

type TokenAuthorizer interface {
	Sign(UserAccount) (string, error)
	Verify(string) error
}

type UserSessionAccessor interface {
	TokenAuthorizer
	Delete(string) error
}
