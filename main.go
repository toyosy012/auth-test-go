package main

import (
	"auth-test/infra"
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
	router := gin.Default()
	router.GET("/users/:id", userAccountController.Get)
	router.GET("/users", userAccountController.List)
	router.POST("/users/new", userAccountController.Create)
	router.PATCH("/users/:id", userAccountController.Update)
	router.DELETE("/users/:id", userAccountController.Delete)
	router.Run("localhost:8080")
}
