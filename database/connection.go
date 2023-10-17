package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DSN = "root:@tcp(127.0.0.1:3306)/meetings?charset=utf8mb4&parseTime=True&loc=Local"
var DB *gorm.DB

func DBconnection() {
	var err error
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	} else {
		log.Println("Connected to Database...")
	}
}

func MigrateDB(t *gorm.Model) {
	DB.AutoMigrate(t)
	log.Println("Database Migration Completed...")
}
