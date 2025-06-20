package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"

	"good-guys/backend/models"
)

var DB *gorm.DB

func Connect() error {

	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("Failed to create data directory", err)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open("data/good-guys.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	DB.AutoMigrate(
		&models.Event{},
		&models.Source{},
		&models.Media{},
	)

	seedDB(DB)
	return nil
}
