package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAccount struct {
	ID    uuid.UUID `gorm:"type:varbinary(36);primaryKey;not null"`
	Email string    `gorm:"unique;not null"`
	Name  string    `gorm:"not null"`
	Hash  string    `gorm:"not null"`
	gorm.Model
}
