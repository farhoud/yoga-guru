package controllers

import (
	"net/http"
	"strconv"
	"time"
	"yoga-guru/internal/models"
	"yoga-guru/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CourseHandler provides methods for course management.
type CourseHandler struct {
	DB *gorm.DB
}

// NewCourseHandler creates a new CourseHandler instance.
func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{DB: db}
}

// CreateCourseRequest defines the request body for creating a course.
type CreateCourseRequest struct {
	Title      string
	CourseType string
	Schedules  []CourseSchedule
	Level      models.CourseLevel
	Price      float64
	Capacity   int
}

type CourseSchedule struct {
	DayOfWeekMask models.DayOfWeekMask
	Recurrence    models.ScheduleRecurrence
	StartTime     utils.CustomTime
	EndTime       utils.CustomTime
}

// CreateCourse godoc
// @Summary Create a new course (Instructor/Admin only)
// @Description Create a new yoga course with details like title, type, schedule, level, price, and capacity.
// @Tags Courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param course body CreateCourseRequest true "Course details"
// @Success 201 {object} models.Course
// @Failure 400 {object} map[string]string "error: Bad request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	instructorID := userIDAny.(uint)

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate CourseLevel
	switch req.Level {
	case models.Beginner, models.Intermediate, models.Advanced:
		// Valid level
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course level specified. Must be 'beginner', 'intermediate', or 'advanced'"})
		return
	}

	schedules := make([]models.Schedule, len(req.Schedules))
	for i, val := range req.Schedules {
		if val.DayOfWeekMask < 0 && val.DayOfWeekMask > 128 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day of week"})
			return
		}

		switch val.Recurrence {
		case models.Weekly, models.BiWeekly, models.MonthlyR:
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Recurrence"})
			return
		}
		schedules[i] = models.Schedule{
			Recurrence: val.Recurrence,
			StartTime:  time.Time(val.StartTime),
			EndTime:    time.Time(val.EndTime),
			DaysMask:   val.DayOfWeekMask,
		}
	}
	course := models.Course{
		Title:        req.Title,
		CourseType:   req.CourseType,
		Schedules:    schedules,
		Level:        req.Level,
		Price:        req.Price,
		Capacity:     req.Capacity,
		InstructorID: instructorID,
	}

	if err := h.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetCourses godoc
// @Summary Get all courses
// @Description Retrieve a list of all available yoga courses.
// @Tags Courses
// @Produce json
// @Success 200 {array} models.Course
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /courses [get]
func (h *CourseHandler) GetCourses(c *gin.Context) {
	var courses []models.Course
	// Preload instructor information for each course
	if err := h.DB.Preload("Instructor").Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// GetCourseByID godoc
// @Summary Get a course by ID
// @Description Retrieve details of a specific yoga course by its ID.
// @Tags Courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 404 {object} map[string]string "error: Course not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course models.Course
	if err := h.DB.Preload("Instructor").First(&course, uint(courseID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}
	c.JSON(http.StatusOK, course)
}

// UpdateCourseRequest defines the request body for updating a course.
type UpdateCourseRequest struct {
	Title      *string
	CourseType *string
	Schedules  []CourseSchedule
	Level      *models.CourseLevel
	Price      *float64
	Capacity   *int
}

// UpdateCourse godoc
// @Summary Update an existing course (Instructor/Admin only)
// @Description Update the details of an existing yoga course. Only the course instructor or an admin can update a course.
// @Tags Courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param course body UpdateCourseRequest true "Updated course details"
// @Success 200 {object} models.Course
// @Failure 400 {object} map[string]string "error: Bad request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: Course not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	currentUserID := userIDAny.(uint)
	userRoleAny := c.MustGet("userRole")
	currentUserRole := userRoleAny.(models.UserRole)

	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var existingCourse models.Course
	if err := h.DB.First(&existingCourse, uint(courseID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}

	// Check if the current user is the instructor of the course or an admin
	if existingCourse.InstructorID != currentUserID && currentUserRole != models.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this course"})
		return
	}

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != nil {
		existingCourse.Title = *req.Title
	}
	if req.CourseType != nil {
		existingCourse.CourseType = *req.CourseType
	}
	schedules := make([]models.Schedule, len(req.Schedules))
	if req.Schedules != nil {
		for i, val := range req.Schedules {
			if val.DayOfWeekMask < 0 && val.DayOfWeekMask > 128 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day of week"})
				return
			}

			switch val.Recurrence {
			case models.Weekly, models.BiWeekly, models.MonthlyR:
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Recurrence"})
				return
			}
			schedules[i] = models.Schedule{
				Recurrence: val.Recurrence,
				StartTime:  time.Time(val.StartTime),
				EndTime:    time.Time(val.EndTime),
				DaysMask:   val.DayOfWeekMask,
			}
		}
	}
	if req.Level != nil {
		// Validate CourseLevel if provided
		switch *req.Level {
		case models.Beginner, models.Intermediate, models.Advanced:
			existingCourse.Level = *req.Level
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course level specified. Must be 'beginner', 'intermediate', or 'advanced'"})
			return
		}
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
			return
		}
		existingCourse.Price = *req.Price
	}
	if req.Capacity != nil {
		if *req.Capacity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Capacity must be greater than 0"})
			return
		}
		existingCourse.Capacity = *req.Capacity
	}

	if err := h.DB.Save(&existingCourse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	c.JSON(http.StatusOK, existingCourse)
}

// DeleteCourse godoc
// @Summary Delete a course (Instructor/Admin only)
// @Description Delete an existing yoga course. Only the course instructor or an admin can delete a course.
// @Tags Courses
// @Security BearerAuth
// @Produce json
// @Param id path int true "Course ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: Course not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	currentUserID := userIDAny.(uint)
	userRoleAny := c.MustGet("userRole")
	currentUserRole := userRoleAny.(models.UserRole)

	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var existingCourse models.Course
	if err := h.DB.First(&existingCourse, uint(courseID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course"})
		return
	}

	// Check if the current user is the instructor of the course or an admin
	if existingCourse.InstructorID != currentUserID && currentUserRole != models.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this course"})
		return
	}

	// Delete associated enrollments first to maintain referential integrity if not set up with CASCADE DELETE
	if err := h.DB.Where("course_id = ?", courseID).Delete(&models.Enrollment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated enrollments"})
		return
	}

	if err := h.DB.Delete(&existingCourse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.Status(http.StatusNoContent)
}
