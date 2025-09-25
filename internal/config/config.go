package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configurations
type Config struct {
	DBPath     string
	Port       string
	JWTSecret  string
}

// LoadConfig reads configuration from environment variables or .env file
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found, using environment variables: %v", err)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./yoga.db" // Default SQLite database path
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set. This is required for authentication.")
	}

	return &Config{
		DBPath:     dbPath,
		Port:       port,
		JWTSecret:  jwtSecret,
	}
}

// You can create a .env file in the root of your project like this:
// DB_PATH=./yoga.db
// PORT=8080
// JWT_SECRET=your_super_secret_jwt_key
