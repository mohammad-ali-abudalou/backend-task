package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Represents The Main User Model.
//
// @Description User Entity That Is Stored In The Database.
type User struct {

	// User Unique Identifier ( UUID ).
	// @Required
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"type:uuid;primaryKey"`

	// Full Name Of The User.
	// @Required
	Name string `json:"name" example:"John Doe" gorm:"not null;size:255" binding:"required"`

	// Email Address ( Unique, Valid format ).
	// @Required
	Email string `json:"email" example:"john.doe@example.com" gorm:"not null;uniqueIndex;size:320" binding:"required,email"`

	// Date Of Birth In YYYY-MM-DD Format ( Must Be In The Past ).
	// @Required
	DateOfBirth time.Time `json:"date_of_birth" example:"1990-05-15" gorm:"type:date;not null" binding:"required"`

	// Group Assignment ( Computed, Read-Only ).
	Group string `json:"group" example:"adult-1" gorm:"not null;index;size:64" readonly:"true"`

	// Timestamp When The Record Was Created.
	CreatedAt time.Time `json:"created_at" example:"2025-09-01T12:00:00Z"`

	// Timestamp When The Record Was Last Updated.
	UpdatedAt time.Time `json:"updated_at" example:"2025-09-01T12:30:00Z"`
}

// Before Create Ensures UUID Is Set Automatically :
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.ID == uuid.Nil {

		u.ID = uuid.New()
	}

	return nil
}
