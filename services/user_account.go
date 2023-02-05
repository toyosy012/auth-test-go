package services

import (
	"auth-test/models"
)

func NewUserAccount(repo models.UserAccountAccessor) UserAccount { return UserAccount{repo: repo} }

type UserAccount struct {
	repo models.UserAccountAccessor
}

func (a UserAccount) Find(id string) (*models.UserAccount, error) {
	account, err := a.repo.Find(id)
	if err != nil {
		return nil, NewApplicationErr(FailedShowUser, err)
	}
	return account, nil
}

func (a UserAccount) FindByEmail(email string) (*models.UserAccount, error) {
	account, err := a.repo.FindByEmail(email)
	if err != nil {
		return nil, NewApplicationErr(FailedShowUser, err)
	}
	return account, nil
}

func (a UserAccount) List() ([]models.UserAccount, error) {
	accounts, err := a.repo.List()
	if err != nil {
		return nil, NewApplicationErr(FailedListUser, err)
	}
	return accounts, nil
}

func (a UserAccount) Create(account models.UserAccount) (*models.UserAccount, error) {
	user, err := a.repo.Insert(account.ID(), account.Email(), account.Name(), account.Password())
	if err != nil {
		return nil, NewApplicationErr(FailedCreateUser, err)
	}
	return user, nil
}

func (a UserAccount) Update(account models.UserAccount) (*models.UserAccount, error) {
	updated, err := a.repo.Update(account)
	if err != nil {
		return nil, NewApplicationErr(FailedUpdateUser, err)
	}
	return updated, nil
}

func (a UserAccount) Delete(id string) error {
	if err := a.repo.Delete(id); err != nil {
		return NewApplicationErr(FailedDeleteUser, err)
	}
	return nil
}
