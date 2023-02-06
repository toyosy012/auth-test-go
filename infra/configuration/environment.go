package configuration

import "time"

type Environment struct {
	Email             string        `envconfig:"EMAIL" required:"true"`
	UserPassword      string        `envconfig:"USER_PASSWORD" required:"true"`
	UserName          string        `envconfig:"USER_NAME" required:"true"`
	User              string        `default:"root"`
	Password          string        `required:"true"`
	Host              string        `default:"0.0.0.0"`
	Port              int           `default:"3306"`
	Name              string        `default:"auth_test"`
	EncryptSecret     string        `envconfig:"ENCRYPT_SECRET" required:"true"`
	RefreshExpiration time.Duration `default:"1h"`
	AccessExpiration  time.Duration `default:"10m"`
	SessionExpiration time.Duration `default:"1h"`
}
