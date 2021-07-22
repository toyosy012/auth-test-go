package models

import "github.com/google/uuid"

type UserAccount struct {
	ID    uuid.UUID
	Email string
	Name  string
	Hash  string
}
