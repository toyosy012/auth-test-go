package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"auth-test/models"
)

type UserAccounts struct {
	ID    string `gorm:"type:varchar(36);primaryKey;not null"`
	Email string `gorm:"unique;not null"`
	Name  string `gorm:"not null"`
	Hash  string `gorm:"not null"`
	gorm.Model
}

func NewUserAccount(id, email, name, passwordHash string) *UserAccounts {
	return &UserAccounts{
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
	var account UserAccounts
	result := r.mysql.Where("id=?", id).First(&account)
	if err := result.Error; err != nil {
		return nil, err
	}

	response := models.NewUserAccount(account.ID, account.Email, account.Name, account.Hash)
	return &response, nil
}

func (r *UserAccountRepository) FindByEmail(email string) (*models.UserAccount, error) {
	var account UserAccounts
	result := r.mysql.Where("email=?", email).First(&account)
	if err := result.Error; err != nil {
		return nil, err
	}

	response := models.NewUserAccount(account.ID, account.Email, account.Name, account.Hash)
	return &response, nil
}

func (r *UserAccountRepository) List() ([]models.UserAccount, error) {
	var accounts []UserAccounts
	result := r.mysql.Find(&accounts)
	if err := result.Error; err != nil {
		return nil, err
	}
	var results []models.UserAccount
	for _, account := range accounts {
		a := models.NewUserAccount(account.ID, account.Email, account.Name, account.Hash)
		results = append(results, a)
	}

	return results, nil
}

func (r *UserAccountRepository) Insert(id, email, name, password string) (*models.UserAccount, error) {
	encryptedPass, err := models.NewEncryption(password)
	if err != nil {
		return nil, err
	}
	account := NewUserAccount(id, email, name, encryptedPass.Hash())

	result := r.mysql.Create(account)
	if err = result.Error; err != nil {
		return nil, err
	}

	var a UserAccounts
	result = r.mysql.Where("email = ?", email).Find(&a)
	if result.Error != nil {
		return nil, result.Error
	}

	response := models.NewUserAccount(a.ID, a.Email, a.Name, a.Hash)
	return &response, nil
}

func (r *UserAccountRepository) Update(account models.UserAccount) (*models.UserAccount, error) {
	encryptedPass, err := models.NewEncryption(account.Password())
	if err != nil {
		return nil, err
	}

	newAccount := UserAccounts{Email: account.Email(), Name: account.Name(), Hash: encryptedPass.Hash()}
	result := r.mysql.Table("user_accounts").Where("id=?", account.ID()).UpdateColumns(newAccount)
	if err = result.Error; err != nil {
		return nil, err
	}

	var a UserAccounts
	result = r.mysql.Where("email = ?", account.Email()).Find(&a)
	if result.Error != nil {
		return nil, result.Error
	}

	response := models.NewUserAccount(a.ID, a.Email, a.Name, a.Hash)
	return &response, nil
}

func (r *UserAccountRepository) Delete(id string) error {
	deletedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	var user models.UserAccount
	result := r.mysql.Where("id = ?", deletedUUID).Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
