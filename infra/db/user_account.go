package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
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

type UserAccountRepositoryImpl struct {
	mysql *gorm.DB
}

func NewUserAccountRepositoryImpl(dsn string) (*UserAccountRepositoryImpl, error) {
	mysqlClient, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(120 * time.Second)

	return &UserAccountRepositoryImpl{
		mysql: mysqlClient,
	}, nil
}

func (i *UserAccountRepositoryImpl) Find(id string) (*models.UserAccount, error) {
	var account UserAccount
	result := i.mysql.Where("id=?", id).First(&account)
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

func (i *UserAccountRepositoryImpl) List() ([]models.UserAccount, error) {
	var accounts []UserAccount
	result := i.mysql.Find(&accounts)
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

func (i *UserAccountRepositoryImpl) Insert(email, name, password string) (*models.UserAccount, error) {
	encryptedPass, err := models.NewEncryptedPassword(password)
	if err != nil {
		return nil, err
	}
	account := NewUserAccount(email, name, encryptedPass.Hash)

	result := i.mysql.Create(account)
	if err := result.Error; err != nil {
		return nil, err
	}

	var a UserAccount
	result = i.mysql.Where("email = ?", account.Email).Find(&a)
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

func (i *UserAccountRepositoryImpl) Update(account models.UserAccount) (*models.UserAccount, error) {
	encryptedPass, err := models.NewEncryptedPassword(account.Password)
	if err != nil {
		return nil, err
	}
	newAccount := UserAccount{Email: account.Email, Name: account.Name, Hash: encryptedPass.Hash}
	result := i.mysql.Table("user_accounts").Where("id=?", account.ID).UpdateColumns(newAccount)
	if err := result.Error; err != nil {
		return nil, err
	}

	var a UserAccount
	result = i.mysql.Where("email = ?", account.Email).Find(&a)
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

func (i *UserAccountRepositoryImpl) Delete(id string) error {
	deletedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil
	}

	user := models.UserAccount{}
	result := i.mysql.Delete(&user, deletedUUID)
	return result.Error
}
