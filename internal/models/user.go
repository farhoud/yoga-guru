package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole defines the possible roles for a user.
type UserRole string

const (
	Admin      UserRole = "admin"
	Instructor UserRole = "instructor"
	Student    UserRole = "student"
)

type UserGender string

const (
	Male   UserGender = "male"
	Female UserGender = "female"
	None   UserGender = ""
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Phone        string         `gorm:"uniqueIndex"`
	PasswordHash string         // Password hash, never expose in JSON
	Role         UserRole
	Profile      Profile
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type Profile struct {
	ID        uuid.UUID `gorm:"primarykey"`
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Gender    UserGender
	AvatarURL string
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
