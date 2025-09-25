package model

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

// Enrollment represents a student's enrollment in a course or package.
type Enrollment struct {
	gorm.Model
	UserID          uint
	User            User
	CourseID        uint // This is the main course this enrollment is for
	Course          Course
	EnrollmentType  EnrollmentType
	StartDate       time.Time
	ExpirationDate  time.Time
	PricePaid       float64
	DiscountApplied float64
	TotalSessions   int // Only for fixed session packages
	SessionsUsed    int // Counter for fixed session packages
	// A user can have many attendance records under this enrollment.
	Attendances []Attendance `gorm:"foreignKey:EnrollmentID"`
	Payments    []Payment    `gorm:"foreignKey:EnrollmentID"`
}

// Attendance tracks whether a user attended a specific course session.
type Attendance struct {
	gorm.Model
	UserID          uint
	User            User
	CourseSessionID uint
	CourseSession   CourseSession
	EnrollmentID    uint
	Enrollment      Enrollment
	Attended        bool // `true` if the user attended, `false` otherwise
	// You might add a field like `NoShow` or `CanceledByStudent` for more detailed tracking.
	RecordedAt time.Time
}

// PaymentStatus defines the status of a payment.
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

// PaymentMethod defines the method of payment.
type PaymentMethod string

const (
	Card          PaymentMethod = "card"
	Cash          PaymentMethod = "cash"
	BankTransfer  PaymentMethod = "bank_transfer"
	OnlinePayment PaymentMethod = "online_payment"
)

// Payment represents a single financial transaction.
type Payment struct {
	gorm.Model
	EnrollmentID  uint // Foreign key to the enrollment this payment is for
	Enrollment    Enrollment
	Amount        float64       // The amount of this specific payment
	Status        PaymentStatus // e.g., 'succeeded', 'failed', 'pending'
	Method        PaymentMethod // e.g., 'card', 'cash'
	TransactionID string        `gorm:"uniqueIndex"` // External ID from payment gateway (e.g., Stripe)
	PaymentDate   time.Time
}
