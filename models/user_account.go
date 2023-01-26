package models

type UserAccount struct {
	ID       string
	Email    string
	Name     string
	Password string
}

type UserAccountRepository interface {
	Find(string) (*UserAccount, error)
	FindByEmail(string) (*UserAccount, error)
	List() ([]UserAccount, error)
	Insert(string, string, string) (*UserAccount, error)
	Update(UserAccount) (*UserAccount, error)
	Delete(string) error
}
