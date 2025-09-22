package models

import (
	"time"

	"gorm.io/gorm"
)

// DayOfWeekMask represents a bitmask for the days of the week.
// Each day is a power of 2.
type DayOfWeekMask int

const (
	Saturday  DayOfWeekMask = 1 << iota // 1 (0000001)
	Sunday                              // 2 (0000010)
	Monday                              // 4 (0000100)
	Tuesday                             // 8 (0001000)
	Wednesday                           // 16 (0010000)
	Thursday                            // 32 (0100000)
	Friday                              // 64 (1000000)
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
	Title        string
	CourseType   string // e.g., Hatha, Vinyasa, Ashtanga
	Level        CourseLevel
	Price        float64    // Price per single session
	Capacity     int        // Max number of students
	InstructorID uint       // ID of the instructor creating the course
	Instructor   User       // GORM association
	Schedules    []Schedule `gorm:"foreignKey:CourseID"`
}

// Schedule defines a specific time, days, and recurrence for a course session.
type Schedule struct {
	gorm.Model
	// Use a single integer field to represent multiple days of the week.
	// Example: A course on Saturday and Sunday would have DaysMask = 3 (1+2).
	DaysMask   DayOfWeekMask
	StartTime  time.Time `gorm:"type:time"`
	EndTime    time.Time `gorm:"type:time"`
	Recurrence string    // e.g., "weekly", "bi-weekly", "monthly"
	CourseID   uint      // Foreign key for the Course
}

// CourseSession represents a single, specific class instance.
// e.g., "Hatha Yoga" on "Monday, October 26, 2025 at 10:00 AM".
type CourseSession struct {
	gorm.Model
	CourseID    uint
	Course      Course
	ScheduledAt time.Time
	// Optional: You could add fields like `Room`, `InstructorID` here if they vary.
	// You might also add a boolean for `IsCancelled`.
	IsCanceled bool
}
