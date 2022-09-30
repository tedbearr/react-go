package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tedbearr/react-go/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {
	errEnv := godotenv.Load()

	if errEnv != nil {
		panic("Failed to laod env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=UTF8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to create connection to database")
	}

	db.AutoMigrate(&entity.User{}, &entity.Posts{})
	return db
}

func CloseDbConnection(db *gorm.DB) {
	dbSql, err := db.DB()

	if err != nil {
		panic("Failed to close connection")
	}

	dbSql.Close()
}
