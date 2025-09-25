package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"yoga-guru/internal/models"
	"yoga-guru/internal/utils"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	Getgorm() *gorm.DB
}

type service struct {
	db  *sql.DB
	orm *gorm.DB
}

var (
	dburl      = os.Getenv("APP_DB_URL")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	gormdb, err := gorm.Open(sqlite.Open(dburl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db, err := gormdb.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	dbInstance = &service{
		db:  db,
		orm: gormdb,
	}

	migrate(gormdb)
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}
	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}

func (s *service) Getgorm() *gorm.DB {
	return s.orm
}

func migrate(db *gorm.DB) *gorm.DB {
	// Auto-migrate the models
	err := db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Course{}, &models.Enrollment{})
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
