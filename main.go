package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"

	"auth-test/infra"
	"auth-test/infra/auth"
	"auth-test/infra/controller"
	"auth-test/infra/db"
	"auth-test/services"
)

func main() {
	var env infra.Environment
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("環境変数の取得に失敗 : %s\n", err.Error())
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
	userAccountController := controller.UserAccountHandler{
		UserAccountService: userAccountSvc,
	}

	tokenAuth := auth.NewTokenAuthentication(env.EncryptSecret, env.AvailabilityTime)
	authSvc := services.NewAuthorizer(userAccountRepo, tokenAuth)

	jwtAuth := controller.NewJWTAuth(authSvc)

	router := gin.Default()
	usersRouter := router.Group("users")
	{
		usersRouter.GET(":id", userAccountController.Get)
		usersRouter.POST("new", userAccountController.Create)
	}
	userSessionRepo := db.NewUserSessionRepo(*dbClient)
	storedAuthSvc := services.NewStoredAuthorization(userAccountRepo, userSessionRepo, env.AvailabilityTime)
	storedAuth := controller.NewStoredAuth(storedAuthSvc)
	v0 := router.Group("v0")
	{
		v0.POST("login", storedAuth.Login)

		v0UsersRouter := v0.Group("users")
		{
			v0UsersRouter.Use(storedAuth.CheckAuthentication).GET("", userAccountController.List)
			authedOwnerRouter := v0UsersRouter.Use(storedAuth.CheckAuthenticatedOwner)
			{
				authedOwnerRouter.PATCH(":id", userAccountController.Update)
				authedOwnerRouter.Use(storedAuth.CheckAuthenticatedOwner).
					DELETE(":id", userAccountController.Delete)
				authedOwnerRouter.Use(storedAuth.CheckAuthenticatedOwner).
					DELETE("logout/:id", storedAuth.Logout)
			}
		}
	}

	v1 := router.Group("v1")
	{
		v1.POST("login", jwtAuth.Login)

		authRouter := v1.Group("users").Use(jwtAuth.CheckAuthentication)
		{
			authRouter.PATCH(":id", userAccountController.Update)
			authRouter.DELETE(":id", userAccountController.Delete)
		}
	}
	router.Run("0.0.0.0:8080")
}
