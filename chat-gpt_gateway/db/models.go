package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID       uint
	Username string `gorm:"unique"`
}

type Chat struct {
	gorm.Model

	ID     string `gorm:"unique"`
	Name   string
	UserID uint
	User   User
}

type VerifiedToken struct {
	gorm.Model

	ID     uint   `gorm:"unique"`
	Token  string `gorm:"unique"`
	UserID uint
	User   User
}

func AutoMigrateDb() {
	// Миграция таблицы, если она еще не существует
	err := DbClient.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	err = DbClient.AutoMigrate(&Chat{})
	if err != nil {
		panic(err)
	}

	err = DbClient.AutoMigrate(&VerifiedToken{})
	if err != nil {
		panic(err)
	}
}
