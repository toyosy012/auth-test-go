package db

import (
	"time"

	"auth-test/models"

	"gorm.io/gorm"
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
	if result.Error != nil {
		return "", nil
	}

	var t Tokens
	result = r.client.
		Where("id = ? AND user_account_id = ?", refreshToken.Value(), refreshToken.AccountID()).
		Limit(1).
		Find(&t)

	if err := result.Error; err != nil {
		return "", err
	}
	return t.ID, nil
}

func (r TokenRepository) Delete(refreshToken string) error {
	if result := r.client.Table("tokens").Delete(Tokens{ID: refreshToken}); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r TokenRepository) FindOwner(refreshToken string, now time.Time) (*models.TokenOwner, error) {
	var token Tokens
	result := r.client.
		Table("tokens").
		Where("id = ? AND ? < expired_at", refreshToken, now.String()).
		Preload("UserAccount").
		Find(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	response := models.NewTokenOwner(token.UserAccount.ID, token.UserAccount.Email)
	return &response, nil
}
