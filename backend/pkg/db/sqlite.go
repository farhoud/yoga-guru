package db

import (
	"log"
	"yoga-backend/config"
	"yoga-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the SQLite database and performs migrations.
func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate the models
	err = db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Course{}, &models.Enrollment{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database connection established and models migrated successfully.")
	return db
}
