package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewClient(dsn string) (*gorm.DB, error) {
	mysqlClient, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	db, err := mysqlClient.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(120 * time.Second)

	return mysqlClient, nil
}
