package models

import (
	"gorm.io/gorm"
)

// CourseLevel defines the difficulty levels for courses.
type CourseLevel string

const (
	Beginner     CourseLevel = "beginner"
	Intermediate CourseLevel = "intermediate"
	Advanced     CourseLevel = "advanced"
)

// Course represents a yoga session or course.
type Course struct {
	gorm.Model
	Title      string      
	CourseType string       // e.g., Hatha, Vinyasa, Ashtanga
	Schedule   string          // e.g., "Every Monday 10:00 AM, Wednesday 6:00 PM" (can be JSON string for more complex schedules)
	Level      CourseLevel 
	Price      float64           // Price per single session
	Capacity   int            // Max number of students
	InstructorID uint       // ID of the instructor creating the course
	Instructor User            // GORM association
}
