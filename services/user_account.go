package services

import "auth-test/models"

type UserAccount struct {
	Repo models.UserAccountRepository
}

func (a UserAccount) Find(id string) (*models.UserAccount, error) {
	account, err := a.Repo.Find(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a UserAccount) List() ([]models.UserAccount, error) {
	accounts, err := a.Repo.List()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (a UserAccount) Create(account models.UserAccount) (*models.UserAccount, error) {
	updated, err := a.Repo.Insert(account.Email, account.Name, account.Password)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (a UserAccount) Update(account models.UserAccount) (*models.UserAccount, error) {
	updated, err := a.Repo.Update(account)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (a UserAccount) Delete(id string) error {
	err := a.Repo.Delete(id)
	return err
}
