package models

import "golang.org/x/crypto/bcrypt"

const (
	EncryptCost = 10
)

type EncryptedPassword struct {
	Hash string
}

func NewEncryptedPassword(password string) (*EncryptedPassword, error) {
	hash, err := hashAndStretch(password)
	if err != nil {
		return nil, err
	}

	return &EncryptedPassword{
		Hash: hash,
	}, nil
}

func (p EncryptedPassword) MatchWith(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(password))
}

func hashAndStretch(password string) (string, error) {
	passwd, err := bcrypt.GenerateFromPassword([]byte(password), EncryptCost)
	if err != nil {
		return "", err
	}

	return string(passwd), nil
}
