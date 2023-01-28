package models

type UserSessionAccessor interface {
	Register(string, string) (string, error)
	Verify(string) error
	Delete(string) error
}
