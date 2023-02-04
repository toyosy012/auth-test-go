package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"

	"auth-test/infra"
	"auth-test/infra/auth"
	"auth-test/infra/controller"
	"auth-test/infra/db"
	"auth-test/services"
)

const (
	PasswordTag = "nist_sp_800_63"
)

func main() {
	var env infra.Environment
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("環境変数の取得に失敗 : %s\n", err.Error())
	}
	validate := validator.New()
	if err != nil {
		log.Fatalf("環境変数の取得に失敗 : %s\n", err.Error())
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err = v.RegisterValidation(PasswordTag, controller.ValidatePassword); err != nil {
			log.Fatalf("パスワードのカスタムバリデーション設定に失敗 : %s\n", err.Error())
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC", env.User, env.Password, env.Host, env.Port, env.Name)
	dbClient, err := db.NewClient(dsn)
	if err != nil {
		log.Fatalf("データベースクライアントの生成に失敗 : %s\n", err.Error())
	}
	userAccountRepo, err := db.NewUserAccountRepository(*dbClient)
	userAccountSvc := services.UserAccount{
		Repo: userAccountRepo,
	}
	userAccountController := controller.NewUserAccountHandler(userAccountSvc, *validate)

	tokenAuth := auth.NewTokenAuthentication(env.EncryptSecret)
	tokenRepo := db.NewTokenRepository(*dbClient)
	tokenAuthSvc := services.NewTokenAuthorization(
		tokenAuth, tokenRepo, userAccountRepo,
		env.RefreshExpiration, env.AccessExpiration,
	)
	tokenAuthController := controller.NewTokenHandler(tokenAuthSvc)

	userSessionRepo := db.NewUserSessionRepo(*dbClient)
	storedAuthSvc := services.NewStoredAuthorization(userAccountRepo, userSessionRepo, env.AvailabilityTime)
	storedAuth := controller.NewStoredAuth(storedAuthSvc)

	router := gin.Default()
	v1 := router.Group("v1")
	{
		v1.POST("login", storedAuth.Login)
		ownedRouter := v1.Group("auth/:id").Use(storedAuth.CheckAuthenticatedOwner)
		{
			ownedRouter.POST("refresh", jwtAuth.Login)
			ownedRouter.DELETE("logout", storedAuth.Logout)
		}

		usersRouter := v1.Group("users")
		{
			usersRouter.POST("new", userAccountController.Create)
			authRouter := usersRouter.Use(jwtAuth.CheckAuthentication)
			{
				authRouter.GET(":id", userAccountController.Get)
				authRouter.GET("", userAccountController.List)
				authRouter.PATCH(":id", userAccountController.Update)
				authRouter.DELETE(":id", userAccountController.Delete)
			}
		}
	}

	router.Run("0.0.0.0:8080")
}
