package main

import (
	"auth-test/infra"
	"auth-test/infra/auth"
	"auth-test/infra/controller"
	"auth-test/infra/db"
	"auth-test/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
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
	userAccountController := controller.UserAccount{
		Service: userAccountSvc,
	}

	authorizer := auth.NewTokenAuthorizer(env.EncryptSecret)
	authSvc := services.NewAuthorizer(userAccountRepo, authorizer)

	loginController := controller.NewLogin(authSvc)

	router := gin.Default()
	router.POST("/login", loginController.Login)
	router.GET("/users/:id", userAccountController.Get)
	router.GET("/users", userAccountController.List)

	authRouter := router.Group("/users").Use(loginController.CheckAuthentication)
	{
		authRouter.POST("new", userAccountController.Create)
		authRouter.PATCH(":id", userAccountController.Update)
		authRouter.DELETE(":id", userAccountController.Delete)
	}
	router.Run("localhost:8080")
}
