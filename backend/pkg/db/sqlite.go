package db

import (
	"log"
	"yoga-backend/config"
	"yoga-backend/models"
	"yoga-backend/pkg/utils"

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

	// Hash the password
	hashedPassword, err := utils.HashPassword("feri1367it")
	if err != nil {
		panic("can not create admin password")
	}
	admin := models.User{
		Phone:        "+989131116127",
		PasswordHash: hashedPassword,
		Role:         models.Admin,
		Profile: models.Profile{
			Name:   "admin",
			Gender: models.Male,
		},
	}

	// Check if user already exists
	var existingUser models.User
	if db.Where("phone = ?", admin.Phone).First(&existingUser).Error == gorm.ErrRecordNotFound {
		if err := db.Create(&admin).Error; err != nil {
			panic("can not create admin user")
		}
	}

	log.Println("Database connection established and models migrated successfully.")
	return db
}
