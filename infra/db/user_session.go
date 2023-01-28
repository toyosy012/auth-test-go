package db

import (
	"time"

	"gorm.io/gorm"

	"auth-test/models"
)

type UserSession struct {
	Owner     string    `gorm:"primaryKey;autoIncrement:false;type:varchar(36);not null"`
	Token     string    `gorm:"primaryKey;autoIncrement:false;not null"`
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
		UserSession{Owner: session.Owner, Token: session.Token, ExpiredAt: session.ExpiredAt},
	)
	if result.Error != nil {
		return "", nil
	}

	var sess UserSession
	result = r.client.
		Where("owner = ? AND token = ?", session.Owner, session.Token).
		Limit(1).
		Find(&sess)

	if err := result.Error; err != nil {
		return "", err
	}

	return sess.Token, nil
}

func (r UserSessionRepository) Verify(token string) error {
	var sess UserSession
	result := r.client.
		Where("token = ? AND ? < expired_at", token, time.Now()).
		First(&sess)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (r UserSessionRepository) FindUser(token string) (string, error) {
	var sess UserSession
	result := r.client.
		Where("token = ? AND ? < expired_at", token, time.Now()).
		First(&sess)
	if err := result.Error; err != nil {
		return "", err
	}

	return sess.Owner, nil
}

func (r UserSessionRepository) Delete(owner, token string) error {
	result := r.client.Delete(UserSession{Owner: owner, Token: token})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
