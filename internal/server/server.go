package server

import (
	"fmt"
	"net/http"
	"time"
	"yoga-guru/docs"
	"yoga-guru/internal/config"
	"yoga-guru/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	cfg *config.Config

	db database.Service
}

func NewServer() *http.Server {
	NewServer := &Server{
		cfg: config.LoadConfig(),
		db:  database.New(),
	}

	// Set up Swagger UI programmatically if not generated
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:" + NewServer.cfg.Port

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", NewServer.cfg.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
