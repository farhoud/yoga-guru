package routes

import (
	"yoga-backend/config"
	"yoga-backend/handlers"
	"yoga-backend/models"
	"yoga-backend/pkg/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gorm.io/gorm"
)

// SetupRouter configures all application routes.
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	courseHandler := handlers.NewCourseHandler(db)
	enrollmentHandler := handlers.NewEnrollmentHandler(db)

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)
	r.GET("/courses", courseHandler.GetCourses)        // Anyone can view courses
	r.GET("/courses/:id", courseHandler.GetCourseByID) // Anyone can view a specific course

	// Authenticated routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(cfg))
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
