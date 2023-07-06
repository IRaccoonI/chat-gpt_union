package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DbClient gorm.DB

func ConnectDb() func() error {

	// БД

	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DbClient = *db

	println("DB connected.")

	// Дополнительные операции с базой данных...

	// Закрываем соединение с базой данных
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	return sqlDB.Close
}
