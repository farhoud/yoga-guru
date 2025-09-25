package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"yoga-guru/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EnrollmentHandler provides methods for enrollment management.
type EnrollmentHandler struct {
	DB *gorm.DB
}

// NewEnrollmentHandler creates a new EnrollmentHandler instance.
func NewEnrollmentHandler(db *gorm.DB) *EnrollmentHandler {
	return &EnrollmentHandler{DB: db}
}

// EnrollRequest defines the request body for course enrollment.
type EnrollRequest struct {
	CourseID       uint
	EnrollmentType models.EnrollmentType
	// Additional fields can be added for specific session dates for 'pre_session' if needed
}

// calculateEnrollmentPrice calculates the total price and applies discounts based on enrollment type.
func calculateEnrollmentPrice(coursePrice float64, enrollmentType models.EnrollmentType) (float64, float64, error) {
	var totalPrice float64
	var discount float64 // Stored as a percentage (e.g., 0.10 for 10%)

	switch enrollmentType {
	case models.PreSession:
		totalPrice = coursePrice // No discount for single session
		discount = 0.0
	case models.Monthly:
		// Assume 4 sessions in a month for calculation, apply a small discount
		totalPrice = coursePrice * 4
		discount = 0.10 // 10% discount for monthly
		totalPrice *= (1 - discount)
	case models.SixMonth:
		// Assume 24 sessions in 6 months
		totalPrice = coursePrice * 24
		discount = 0.20 // 20% discount for 6-month
		totalPrice *= (1 - discount)
	case models.Yearly:
		// Assume 48 sessions in a year
		totalPrice = coursePrice * 48
		discount = 0.30 // 30% discount for yearly
		totalPrice *= (1 - discount)
	default:
		return 0, 0, fmt.Errorf("invalid enrollment type: %s", enrollmentType)
	}

	return totalPrice, discount, nil
}

// EnrollInCourse godoc
// @Summary Enroll a student in a course (Student only)
// @Description Allows a student to enroll in a yoga course with various enrollment packages.
// @Tags Enrollments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param enrollment body EnrollRequest true "Enrollment details"
// @Success 201 {object} models.Enrollment
// @Failure 400 {object} map[string]string "error: Bad request"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: Course not found"
// @Failure 409 {object} map[string]string "error: Already enrolled or course full"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /enrollments [post]
func (h *EnrollmentHandler) EnrollInCourse(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	studentID := userIDAny.(uint)

	var req EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch course details
	var course models.Course
	if err := h.DB.First(&course, req.CourseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch course details"})
		return
	}

	// Check course capacity
	var currentEnrollments int64
	h.DB.Model(&models.Enrollment{}).Where("course_id = ?", req.CourseID).Count(&currentEnrollments)
	if int(currentEnrollments) >= course.Capacity {
		c.JSON(http.StatusConflict, gin.H{"error": "Course is full"})
		return
	}

	// Check if student is already enrolled in this course for the same period (optional, depending on business logic)
	// For simplicity, we'll prevent duplicate enrollments for any type for now.
	var existingEnrollment models.Enrollment
	if h.DB.Where("user_id = ? AND course_id = ?", studentID, req.CourseID).First(&existingEnrollment).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You are already enrolled in this course"})
		return
	}

	// Calculate price and discount
	totalPrice, discount, err := calculateEnrollmentPrice(course.Price, req.EnrollmentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	var endDate time.Time

	switch req.EnrollmentType {
	case models.PreSession:
		endDate = now // For pre-session, end date might just be the session date itself or not applicable
	case models.Monthly:
		endDate = now.AddDate(0, 1, 0) // 1 month from now
	case models.SixMonth:
		endDate = now.AddDate(0, 6, 0) // 6 months from now
	case models.Yearly:
		endDate = now.AddDate(1, 0, 0) // 1 year from now
	}

	enrollment := models.Enrollment{
		UserID:          studentID,
		CourseID:        req.CourseID,
		EnrollmentType:  req.EnrollmentType,
		StartDate:       now,
		ExpirationDate:  endDate,
		PricePaid:       totalPrice,
		DiscountApplied: discount,
	}

	if err := h.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll in course"})
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

// GetStudentEnrollments godoc
// @Summary Get student's enrollments
// @Description Retrieve a list of all courses a student is enrolled in.
// @Tags Enrollments
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Enrollment
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /enrollments/me [get]
func (h *EnrollmentHandler) GetStudentEnrollments(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	studentID := userIDAny.(uint)

	var enrollments []models.Enrollment
	if err := h.DB.Preload("Course.Instructor").Where("user_id = ?", studentID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}

// GetEnrollmentByID godoc
// @Summary Get enrollment by ID
// @Description Retrieve details of a specific enrollment by its ID. (Admin/Enrolled Student only)
// @Tags Enrollments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Enrollment ID"
// @Success 200 {object} models.Enrollment
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: Enrollment not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /enrollments/{id} [get]
func (h *EnrollmentHandler) GetEnrollmentByID(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	currentUserID := userIDAny.(uint)
	userRoleAny := c.MustGet("userRole")
	currentUserRole := userRoleAny.(models.UserRole)

	enrollmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	var enrollment models.Enrollment
	if err := h.DB.Preload("User").Preload("Course.Instructor").First(&enrollment, uint(enrollmentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollment"})
		return
	}

	// Only admin or the enrolled student can view this enrollment
	if currentUserRole != models.Admin && enrollment.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this enrollment"})
		return
	}

	c.JSON(http.StatusOK, enrollment)
}

// CancelEnrollment godoc
// @Summary Cancel an enrollment (Student/Admin only)
// @Description Allows a student to cancel their enrollment, or an admin to cancel any enrollment.
// @Tags Enrollments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Enrollment ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 403 {object} map[string]string "error: Forbidden"
// @Failure 404 {object} map[string]string "error: Enrollment not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /enrollments/{id} [delete]
func (h *EnrollmentHandler) CancelEnrollment(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	currentUserID := userIDAny.(uint)
	userRoleAny := c.MustGet("userRole")
	currentUserRole := userRoleAny.(models.UserRole)

	enrollmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	var enrollment models.Enrollment
	if err := h.DB.First(&enrollment, uint(enrollmentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollment"})
		return
	}

	// Only admin or the enrolled student can cancel this enrollment
	if currentUserRole != models.Admin && enrollment.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to cancel this enrollment"})
		return
	}

	if err := h.DB.Delete(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel enrollment"})
		return
	}

	c.Status(http.StatusNoContent)
}
