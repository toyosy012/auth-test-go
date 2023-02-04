package models

func NewUserAccount(id, email, name, password string) UserAccount {
	return UserAccount{id: id, email: email, name: name, password: password}
}

type UserAccount struct {
	id       string
	email    string
	name     string
	password string
}

func (a UserAccount) ID() string       { return a.id }
func (a UserAccount) Email() string    { return a.email }
func (a UserAccount) Name() string     { return a.name }
func (a UserAccount) Password() string { return a.password }

type UserAccountAccessor interface {
	Find(string) (*UserAccount, error)
	FindByEmail(string) (*UserAccount, error)
	List() ([]UserAccount, error)
	Insert(string, string, string, string) (*UserAccount, error)
	Update(UserAccount) (*UserAccount, error)
	Delete(string) error
}
