package storage

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/yassineidr02/ovh-prep/"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the User model
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully.")
}
