package infra

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/gorm"

	"auth-test/infra/auth"
	"auth-test/infra/configuration"
	"auth-test/infra/controller"
	"auth-test/infra/db"
	"auth-test/services"
)

const (
	PasswordTag = "nist_sp_800_63"
)

func Run() error {
	var env configuration.Environment
	err := envconfig.Process("", &env)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		env.User, env.Password, env.Host, env.Port, env.Name,
	)
	dbClient, err := db.NewClient(dsn)
	if err != nil {
		return err

	}

	validate := validator.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err = v.RegisterValidation(PasswordTag, controller.ValidatePassword); err != nil {
			return err
		}
	}

	router, err := setUpRouter(env, *dbClient, *validate)
	if err != nil {
		return err
	}

	return router.Run("0.0.0.0:8080")
}

func setUpRouter(env configuration.Environment, dbClient gorm.DB, validate validator.Validate) (*gin.Engine, error) {
	userAccountRepo := db.NewUserAccountRepository(dbClient)
	userAccountSvc := services.NewUserAccount(userAccountRepo)
	userAccountController := controller.NewUserAccountHandler(userAccountSvc, validate)

	tokenAuth := auth.NewTokenAuthorization(env.EncryptSecret)
	tokenRepo := db.NewTokenRepository(dbClient)
	tokenAuthSvc := services.NewTokenAuthorization(
		tokenAuth, tokenRepo, userAccountRepo, env.RefreshExpiration, env.AccessExpiration,
	)
	tokenAuthController := controller.NewTokenHandler(tokenAuthSvc)

	userSessionRepo := db.NewUserSessionRepo(dbClient)
	userSessionSvc := services.NewSessionAuthorization(userAccountRepo, userSessionRepo, env.SessionExpiration)
	userSessionController := controller.NewSessionAuth(userSessionSvc)

	router := gin.Default()
	v1 := router.Group("v1")
	usersRouter := v1.Group("users") // デバック用APIのため各認証グループ外に設定
	{
		usersRouter.GET("", userAccountController.List)
		usersRouter.POST("new", userAccountController.Create)
	}

	{
		sessionRouter := v1.Group("session")
		sessionRouter.POST("login", userSessionController.Login)
		sessionRouter.Use(userSessionController.CheckAuthenticatedOwner).DELETE("logout/:id", userSessionController.Logout)
		{
			r := sessionRouter.Group("users").Use(userSessionController.CheckAuthenticatedOwner)
			{
				r.GET(":id", userAccountController.Get)
				r.PUT(":id", userAccountController.Update)
				r.DELETE(":id", userAccountController.Delete)
			}
		}

		authRouter := v1.Group("auth")
		authRouter.POST("claim", tokenAuthController.Claim)
		authRouter.POST("refresh", tokenAuthController.Refresh)
		{
			r := authRouter.Group("users").Use(tokenAuthController.VerifyIDToken)
			{
				r.GET(":id", userAccountController.Get)
				r.PUT(":id", userAccountController.Update)
				r.DELETE(":id", userAccountController.Delete)
			}
		}
	}

	return router, nil
}
