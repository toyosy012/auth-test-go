package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"auth-test/models"
)

type UserAccount struct {
	ID    uuid.UUID `gorm:"type:varbinary(36);primaryKey;not null"`
	Email string    `gorm:"unique;not null"`
	Name  string    `gorm:"not null"`
	Hash  string    `gorm:"not null"`
	gorm.Model
}

func NewUserAccount(email, name, passwordHash string) *UserAccount {
	id := uuid.New()
	return &UserAccount{
		ID:    id,
		Email: email,
		Name:  name,
		Hash:  passwordHash,
	}
}

type UserAccountRepository struct {
	mysql gorm.DB
}

func NewUserAccountRepository(client gorm.DB) (*UserAccountRepository, error) {
	return &UserAccountRepository{
		mysql: client,
	}, nil
}

func (r *UserAccountRepository) Find(id string) (*models.UserAccount, error) {
	var account UserAccount
	result := r.mysql.Where("id=?", id).First(&account)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &models.UserAccount{
		ID:       account.ID.String(),
		Email:    account.Email,
		Name:     account.Name,
		Password: account.Hash,
	}, nil
}

func (r *UserAccountRepository) FindByEmail(email string) (*models.UserAccount, error) {
	var account UserAccount
	result := r.mysql.Where("email=?", email).First(&account)
	if err := result.Error; err != nil {
		return nil, err
	}

	return &models.UserAccount{
		ID:       account.ID.String(),
		Email:    account.Email,
		Name:     account.Name,
		Password: account.Hash,
	}, nil
}

func (r *UserAccountRepository) List() ([]models.UserAccount, error) {
	var accounts []UserAccount
	result := r.mysql.Find(&accounts)
	if err := result.Error; err != nil {
		return nil, err
	}
	var results []models.UserAccount
	for _, account := range accounts {
		a := models.UserAccount{
			ID:       account.ID.String(),
			Email:    account.Email,
			Name:     account.Name,
			Password: account.Hash,
		}
		results = append(results, a)
	}

	return results, nil
}

func (r *UserAccountRepository) Insert(email, name, password string) (*models.UserAccount, error) {
	encryptedPass, err := models.NewEncryptedPassword(password)
	if err != nil {
		return nil, err
	}
	account := NewUserAccount(email, name, encryptedPass.Hash)

	result := r.mysql.Create(account)
	if err := result.Error; err != nil {
		return nil, err
	}

	var a UserAccount
	result = r.mysql.Where("email = ?", account.Email).Find(&a)
	if result.Error != nil {
		return nil, result.Error
	}

	return &models.UserAccount{
		ID:       a.ID.String(),
		Email:    a.Email,
		Name:     a.Name,
		Password: a.Hash,
	}, nil
}

func (r *UserAccountRepository) Update(account models.UserAccount) (*models.UserAccount, error) {
	encryptedPass, err := models.NewEncryptedPassword(account.Password)
	if err != nil {
		return nil, err
	}
	newAccount := UserAccount{Email: account.Email, Name: account.Name, Hash: encryptedPass.Hash}
	result := r.mysql.Table("user_accounts").Where("id=?", account.ID).UpdateColumns(newAccount)
	if err := result.Error; err != nil {
		return nil, err
	}

	var a UserAccount
	result = r.mysql.Where("email = ?", account.Email).Find(&a)
	if result.Error != nil {
		return nil, result.Error
	}

	return &models.UserAccount{
		ID:       a.ID.String(),
		Email:    a.Email,
		Name:     a.Name,
		Password: a.Hash,
	}, nil
}

func (r *UserAccountRepository) Delete(id string) error {
	deletedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil
	}

	user := models.UserAccount{}
	result := r.mysql.Delete(&user, deletedUUID)
	return result.Error
}
