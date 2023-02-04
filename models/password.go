package models

import "golang.org/x/crypto/bcrypt"

const (
	EncryptCost = 10
)

type EncryptedPassword struct {
	hash string
}

func NewEncryptedPassword(hash string) EncryptedPassword { return EncryptedPassword{hash: hash} }

func NewEncryption(password string) (*EncryptedPassword, error) {
	hash, err := hashAndStretch(password)
	if err != nil {
		return nil, err
	}

	return &EncryptedPassword{
		hash: hash,
	}, nil
}

func (p EncryptedPassword) MatchWith(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(password))
}

func (p EncryptedPassword) Hash() string { return p.hash }

func hashAndStretch(password string) (string, error) {
	passwd, err := bcrypt.GenerateFromPassword([]byte(password), EncryptCost)
	if err != nil {
		return "", err
	}

	return string(passwd), nil
}
