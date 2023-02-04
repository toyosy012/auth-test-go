package db

import (
	"time"

	"gorm.io/gorm"

	"auth-test/models"
)

type UserSessions struct {
	ID        string    `gorm:"type:varchar(36);primaryKey;not null"`
	UserID    string    `gorm:"type:varchar(36);type:varchar(36);not null"`
	ExpiredAt time.Time `gorm:"type:datetime(0);not null"`
	CreatedAt time.Time `gorm:"type:datetime(0);not null;default:current_timestamp"`
}

type Token struct {
	Token string
}

func NewUserSessionRepo(client gorm.DB) UserSessionRepository {
	return UserSessionRepository{
		client: client,
	}
}

type UserSessionRepository struct {
	client gorm.DB
}

func (r UserSessionRepository) Register(session models.Session) (string, error) {
	result := r.client.Create(
		UserSessions{
			UserID:    session.Owner(),
			ID:        session.Token(),
			ExpiredAt: session.ExpiredAt(),
		},
	)
	if result.Error != nil {
		return "", nil
	}

	var sess UserSessions
	result = r.client.
		Where("id = ? AND user_id = ?", session.Token(), session.Owner()).
		Limit(1).
		Find(&sess)

	if err := result.Error; err != nil {
		return "", err
	}

	return sess.ID, nil
}

func (r UserSessionRepository) Verify(token string) error {
	var sess UserSessions
	result := r.client.
		Where("token = ? AND ? < expired_at", token, time.Now()).
		First(&sess)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (r UserSessionRepository) FindOwner(token string) (string, error) {
	var sess UserSessions
	result := r.client.
		Where("id = ? AND ? < expired_at", token, time.Now()).
		First(&sess)
	if err := result.Error; err != nil {
		return "", err
	}

	return sess.UserID, nil
}

func (r UserSessionRepository) Delete(owner, token string) error {
	result := r.client.Delete(UserSessions{
		UserID: owner,
		ID:     token,
	})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
