package infra

type Environment struct {
	Email         string `envconfig:"EMAIL" required:"true"`
	UserPassword  string `envconfig:"USER_PASSWORD" required:"true"`
	UserName      string `envconfig:"USER_NAME" required:"true"`
	User          string `default:"root"`
	Password      string `default:"mysql"`
	Host          string `default:"0.0.0.0"`
	Port          int    `default:"3306"`
	Name          string `default:"auth_test"`
	EncryptSecret string `default:"secret"`
}
