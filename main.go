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

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		env.User, env.Password, env.Host, env.Port, env.Name,
	)
	dbClient, err := db.NewClient(dsn)
	if err != nil {
		log.Fatalf("データベースクライアントの生成に失敗 : %s\n", err.Error())
	}
	userAccountRepo, err := db.NewUserAccountRepository(*dbClient)
	userAccountSvc := services.NewUserAccount(userAccountRepo)
	userAccountController := controller.NewUserAccountHandler(userAccountSvc, *validate)

	tokenAuth := auth.NewTokenAuthorization(env.EncryptSecret)
	tokenRepo := db.NewTokenRepository(*dbClient)
	tokenAuthSvc := services.NewTokenAuthorization(
		tokenAuth, tokenRepo, userAccountRepo,
		env.RefreshExpiration, env.AccessExpiration,
	)
	tokenAuthController := controller.NewTokenHandler(tokenAuthSvc)

	userSessionRepo := db.NewUserSessionRepo(*dbClient)
	userSessionSvc := services.NewSessionAuthorization(userAccountRepo, userSessionRepo, env.SessionExpiration)
	userSessionController := controller.NewSessionAuth(userSessionSvc)

	router := gin.Default()
	v1 := router.Group("v1")
	usersRouter := v1.Group("users")
	{
		usersRouter.GET("", userAccountController.List) // デバック用APIのため各認証グループ外ルーティングに設定
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

	router.Run("0.0.0.0:8080")
}
