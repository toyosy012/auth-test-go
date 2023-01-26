package models

import "time"

type TokenAuthorizer interface {
	Sign(UserAccount, time.Time) (string, error)
	Verify(string, time.Time) error
}
