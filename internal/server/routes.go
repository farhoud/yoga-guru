package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"yoga-guru/internal/controllers"
	"yoga-guru/internal/middleware"
	"yoga-guru/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/coder/websocket"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/websocket", s.websocketHandler)
	// Initialize handlers
	authHandler := controllers.NewAuthHandler(s.db.Getgorm(), s.cfg)
	userHandler := controllers.NewUserHandler(s.db.Getgorm())
	courseHandler := controllers.NewCourseHandler(s.db.Getgorm())
	enrollmentHandler := controllers.NewEnrollmentHandler(s.db.Getgorm())

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)
	r.GET("/courses", courseHandler.GetCourses)        // Anyone can view courses
	r.GET("/courses/:id", courseHandler.GetCourseByID) // Anyone can view a specific course

	// Authenticated routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(s.cfg))
	{
		// User routes
		authorized.GET("/users/me", userHandler.GetCurrentUserProfile)

		// Admin-only user routes
		adminGroup := authorized.Group("/")
		adminGroup.Use(middleware.AuthorizeRole(models.Admin))
		{
			adminGroup.PUT("/users/:id/role", userHandler.UpdateUserRole)
		}

		// Instructor and Admin routes for course management
		instructorAdminGroup := authorized.Group("/courses")
		instructorAdminGroup.Use(middleware.AuthorizeRole(models.Instructor, models.Admin))
		{
			instructorAdminGroup.POST("", courseHandler.CreateCourse)
			instructorAdminGroup.PUT("/:id", courseHandler.UpdateCourse)
			instructorAdminGroup.DELETE("/:id", courseHandler.DeleteCourse)
		}

		// Student and Admin routes for enrollments
		studentAdminGroup := authorized.Group("/enrollments")
		studentAdminGroup.Use(middleware.AuthorizeRole(models.Student, models.Admin))
		{
			studentAdminGroup.POST("", enrollmentHandler.EnrollInCourse)
			studentAdminGroup.GET("/me", enrollmentHandler.GetStudentEnrollments)
			studentAdminGroup.GET("/:id", enrollmentHandler.GetEnrollmentByID)
			studentAdminGroup.DELETE("/:id", enrollmentHandler.CancelEnrollment)
		}
	}

	// Swagger documentation route
	// The url point to the new file created by swag init
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) websocketHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	socket, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}
