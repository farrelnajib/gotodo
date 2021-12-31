package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DBNAME")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := "3306"

	dbUri := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName)
	fmt.Print(dbUri)

	connection, err := gorm.Open(mysql.Open(dbUri), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Print(err)
	}

	db = connection
	db.AutoMigrate(&Activity{}, &Todo{})
}

func GetDB() *gorm.DB {
	return db
}
