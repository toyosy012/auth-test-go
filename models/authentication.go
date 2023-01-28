package models

type Authorizer interface {
	Sign(UserAccount) (string, error)
	Verify(string) error
}
