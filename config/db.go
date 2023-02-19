package config

import (
	"github.com/Seals29/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// import "gorm.io/driver/postgres"
// ref: https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open("user=postgres password=123 dbname=newegg port=9999"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.ResetUser{})
	db.AutoMigrate(&models.UserSubscribe{})
	DB = db
}
