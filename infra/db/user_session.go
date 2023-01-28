package db

func NewUserSessionRepo() UserSessionRepository {
	return UserSessionRepository{}
}

type UserSessionRepository struct {
}

func (r UserSessionRepository) Register(string, string) (string, error) {
	return "", nil
}

func (r UserSessionRepository) Verify(string) error {
	return nil
}

func (r UserSessionRepository) Delete(string) error {
	return nil
}
