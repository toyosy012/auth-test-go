package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"auth-test/infra/configuration"
	"auth-test/infra/db"
	"auth-test/models"
)

func main() {
	var env configuration.Environment
	err := envconfig.Process("", &env)

	if err != nil {
		log.Fatalf("環境変数の読み込み失敗。: %s \n", err.Error())
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.User, env.Password, env.Host, env.Port, env.Name)
	mysqlDB, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("DBへの接続に失敗。: %s \n", err.Error())
	}

	err = mysqlDB.AutoMigrate(&db.UserSessions{})
	if err != nil {
		log.Fatalf("テーブルのマイグレーションに失敗。: %s \n", err.Error())
	}

	err = mysqlDB.AutoMigrate(&db.UserAccounts{})
	if err != nil {
		log.Fatalf("テーブルのマイグレーションに失敗。: %s \n", err.Error())
	}

	err = mysqlDB.AutoMigrate(&db.Tokens{})
	if err != nil {
		log.Fatalf("テーブルのマイグレーションに失敗。: %s \n", err.Error())
	}

	newID := uuid.New()
	encrypted, err := models.NewEncryption(env.UserPassword)
	if err != nil {
		log.Fatal("パスワードの保存に失敗")
	}
	mysqlDB.Create(&db.UserAccounts{
		ID:    newID.String(),
		Email: env.Email,
		Name:  env.UserName,
		Hash:  encrypted.Hash(),
	})
}
