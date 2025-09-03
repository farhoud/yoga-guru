package main

import (
	"log"
	"yoga-backend/config"
	"yoga-backend/docs"
	"yoga-backend/pkg/db"
	"yoga-backend/routes"

	_ "gorm.io/gorm" // Blank import for GORM to help 'swag' tool resolve gorm.Model
)

// @title Yoga Backend API
// @version 1.0
// @description This is a backend API for a yoga session management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database := db.InitDB(cfg)

	// Setup Gin router
	router := routes.SetupRouter(database, cfg)

	// Set up Swagger UI programmatically if not generated
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port

	// Run the server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
