package config

import (
	"vendingmachine/domains"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func DBConnect() {

	dsn := "root:r00t@tcp(127.0.0.1:3306)/vending_machine_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	Database = db

	if err != nil {
		panic("failed to connect database")
	}

	runMigrations()
}

func runMigrations() {
	Database.AutoMigrate(&domains.Product{})
	Database.AutoMigrate(&domains.UserDeposit{})
}
