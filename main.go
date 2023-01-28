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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.User, env.Password, env.Host, env.Port, env.Name)
	userAccountRepo, err := db.NewUserAccountRepositoryImpl(dsn)
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
	router.GET("/users/:id", userAccountController.Get)
	router.GET("/users", userAccountController.List)
	router.POST("new", userAccountController.Create)

	userSessionRepo := db.NewUserSessionRepo()
	storedAuthSvc := services.NewStoredAuthorization(userAccountRepo, userSessionRepo)
	storedAuth := controller.NewStoredAuth(storedAuthSvc)
	v0 := router.Group("/v0")
	{
		v0.POST("/login", storedAuth.Login)

		storedAuthRouter := v0.Group("/users").Use(storedAuth.CheckAuthentication)
		{
			storedAuthRouter.PATCH(":id", userAccountController.Update)
			storedAuthRouter.DELETE(":id", userAccountController.Delete)
		}
	}

	v1 := router.Group("/v1")
	{
		v1.POST("/login", jwtAuth.Login)

		authRouter := v1.Group("/users").Use(jwtAuth.CheckAuthentication)
		{
			authRouter.PATCH(":id", userAccountController.Update)
			authRouter.DELETE(":id", userAccountController.Delete)
		}
	}
	router.Run("localhost:8080")
}
