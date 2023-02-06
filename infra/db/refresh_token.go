package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"auth-test/models"
	"auth-test/services"
)

const (
	MySQLDuplicateEntry = 1062
	NoDeleteRecords     = 0
)

type Tokens struct {
	ID            string       `gorm:"type:varchar(36);primaryKey;not null"`
	UserAccountID string       `gorm:"type:varchar(36);not null;constraint:OnDelete:CASCADE"`
	ExpiredAt     time.Time    `gorm:"type:datetime(0);not null"`
	CreatedAt     time.Time    `gorm:"type:datetime(0);not null;default:current_timestamp"`
	UserAccount   UserAccounts `gorm:"foreignKey:UserAccountID;constraint:OnDelete:CASCADE"`
}

func NewTokenRepository(client gorm.DB) TokenRepository {
	return TokenRepository{
		client: client,
	}
}

type TokenRepository struct {
	client gorm.DB
}

func (r TokenRepository) Insert(refreshToken models.RefreshTokenInput) (string, error) {
	result := r.client.
		Create(
			Tokens{
				ID:            refreshToken.Value(),
				UserAccountID: refreshToken.AccountID(),
				ExpiredAt:     refreshToken.ExpiredAt(),
			},
		)
	if err := result.Error; err != nil {
		switch {
		case err.(*mysql.MySQLError).Number == MySQLDuplicateEntry:
			return "", services.NewApplicationErr(services.DuplicateToken, err)
		default:
			return "", services.NewApplicationErr(services.InternalServerErr, err)
		}
	}

	var t Tokens
	result = r.client.
		Where("id = ? AND user_account_id = ?", refreshToken.Value(), refreshToken.AccountID()).
		First(&t)

	if err := result.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return "", services.NewApplicationErr(services.NoTokenRecord, err)
		default:
			return "", services.NewApplicationErr(services.InternalServerErr, err)
		}
	}
	return t.ID, nil
}

func (r TokenRepository) Delete(refreshToken string) error {
	result := r.client.
		Unscoped().
		Table("tokens").
		Delete(Tokens{ID: refreshToken})
	if result.RowsAffected == NoDeleteRecords {
		return services.NewApplicationErr(services.NoTokenRecord, fmt.Errorf("削除対象: %s", refreshToken))
	} else if result.Error != nil {
		return services.NewApplicationErr(services.InternalServerErr, result.Error)
	}
	return nil
}

func (r TokenRepository) FindOwner(refreshToken string, now time.Time) (*models.TokenOwner, error) {
	var token Tokens
	result := r.client.
		Table("tokens").
		Where("id = ? AND ? < expired_at", refreshToken, now.String()).
		Preload("UserAccount").
		First(&token)
	if err := result.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, services.NewApplicationErr(services.NoUserRecord, err)
		default:
			return nil, services.NewApplicationErr(services.InternalServerErr, err)
		}
	}
	response := models.NewTokenOwner(token.UserAccount.ID, token.UserAccount.Email)
	return &response, nil
}
