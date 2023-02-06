package db

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"auth-test/models"
	"auth-test/services"
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

func NewUserAccountRepository(client gorm.DB) *UserAccountRepository {
	return &UserAccountRepository{mysql: client}
}

func (r *UserAccountRepository) Find(id string) (*models.UserAccount, error) {
	var account UserAccounts
	result := r.mysql.Where("id=?", id).First(&account)
	if err := result.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, services.NewApplicationErr(services.NoUserRecord, err)
		default:
			return nil, services.NewApplicationErr(services.InternalServerErr, err)
		}
	}

	response := models.NewUserAccount(account.ID, account.Email, account.Name, account.Hash)
	return &response, nil
}

func (r *UserAccountRepository) FindByEmail(email string) (*models.UserAccount, error) {
	var account UserAccounts
	result := r.mysql.Where("email=?", email).First(&account)
	if err := result.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, services.NewApplicationErr(services.NoUserEmail, err)
		default:
			return nil, services.NewApplicationErr(services.InternalServerErr, err)
		}
	}

	response := models.NewUserAccount(account.ID, account.Email, account.Name, account.Hash)
	return &response, nil
}

func (r *UserAccountRepository) List() ([]models.UserAccount, error) {
	var accounts []UserAccounts
	result := r.mysql.Find(&accounts)
	if err := result.Error; err != nil {
		return nil, services.NewApplicationErr(services.NoUsersRecord, err)
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
		return nil, services.NewApplicationErr(services.TooLongPassword, err)
	}
	account := NewUserAccount(id, email, name, encryptedPass.Hash())

	result := r.mysql.Create(account)
	if err = result.Error; err != nil {
		switch {
		case err.(*mysql.MySQLError).Number == MySQLDuplicateEntry:
			return nil, services.NewApplicationErr(services.DuplicateUserEmail, err)
		default:
			return nil, services.NewApplicationErr(services.InternalServerErr, err)
		}
	}

	var a UserAccounts
	result = r.mysql.Where("email = ?", email).Find(&a)
	if err = result.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, services.NewApplicationErr(services.NoUserRecord, err)
		default:
			return nil, services.NewApplicationErr(services.InternalServerErr, err)
		}
	}

	response := models.NewUserAccount(a.ID, a.Email, a.Name, a.Hash)
	return &response, nil
}

func (r *UserAccountRepository) Update(account models.UserAccount) (*models.UserAccount, error) {
	encryptedPass, err := models.NewEncryption(account.Password())
	if err != nil {
		return nil, services.NewApplicationErr(services.TooLongPassword, err)
	}

	newAccount := UserAccounts{Email: account.Email(), Name: account.Name(), Hash: encryptedPass.Hash()}
	var a UserAccounts
	result := r.mysql.Table("user_accounts").Where("id = ?", account.ID()).UpdateColumns(newAccount).First(&a)
	if err = result.Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, services.NewApplicationErr(services.NoUserRecord, err)
		case err.(*mysql.MySQLError).Number == MySQLDuplicateEntry:
			return nil, services.NewApplicationErr(services.DuplicateUserEmail, err)
		default:
			return nil, services.NewApplicationErr(services.InternalServerErr, err)
		}
	}

	response := models.NewUserAccount(a.ID, a.Email, a.Name, a.Hash)
	return &response, nil
}

func (r *UserAccountRepository) Delete(id string) error {
	deletedUUID, err := uuid.Parse(id)
	if err != nil {
		return services.NewApplicationErr(services.InvalidUUIDFormat, err)
	}

	result := r.mysql.Unscoped().Delete(&UserAccounts{}, deletedUUID)
	if result.RowsAffected == NoDeleteRecords {
		return services.NewApplicationErr(services.NoUserRecord, fmt.Errorf("削除対象ID: %s", id))
	} else if result.Error != nil {
		return services.NewApplicationErr(services.InternalServerErr, result.Error)
	}
	return nil
}
