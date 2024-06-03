package database

import (
	"finalapi/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	database, error := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/bank_user_db?parseTime=true&loc=Local"))
	if error != nil {
		panic(error)
	}
	database.AutoMigrate(&models.Photo{}, &models.User{})

	DB = database
}
