package dataSource

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func connect(host string, user string, pass string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, pass, dbName)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func NewDataSource(host string, user string, pass string, dbName string) error {
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
