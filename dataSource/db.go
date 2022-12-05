package dataSource

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func NewDataSource(host string, user string, pass string, dbName string) error {
	if db != nil {
		return nil
	}

	handle, err := connect(host, user, pass, dbName)

	if err != nil {
		return err
	}

	err = handle.AutoMigrate(&Country{}, &Person{})

	if err != nil {
		return err
	}

	db = handle

	return err
}

func DB() *gorm.DB {
	return db
}

func connect(host string, user string, pass string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, pass, dbName)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
