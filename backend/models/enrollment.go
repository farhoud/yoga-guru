package models

import (
	"time"

	"gorm.io/gorm"
)

// EnrollmentType defines the different types of enrollment packages.
type EnrollmentType string

const (
	PreSession EnrollmentType = "pre_session"
	Monthly    EnrollmentType = "monthly"
	SixMonth   EnrollmentType = "six_month"
	Yearly     EnrollmentType = "yearly"
)

// Enrollment represents a student's enrollment in a course.
type Enrollment struct {
	gorm.Model
	UserID        uint           
	User          User           
	CourseID      uint           
	Course        Course         
	EnrollmentType EnrollmentType 
	StartDate     time.Time      
	EndDate       time.Time       // End date for subscriptions (monthly, yearly etc.)
	PricePaid     float64        
	DiscountApplied float64       // Discount percentage applied
	// For pre-session, you might also want to store specific session dates for 'pre_session' if the schedule isn't fixed
	// e.g., SessionDates []time.Time 
}
