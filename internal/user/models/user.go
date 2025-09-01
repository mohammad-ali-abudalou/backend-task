package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Is The Main User Model :
type User struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Email       string    `json:"email" gorm:"not null;uniqueIndex;size:320"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"type:date;not null"`
	Group       string    `json:"group" gorm:"not null;index;size:64"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Before Create Ensures UUID Is Set Automatically :
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.ID == uuid.Nil {

		u.ID = uuid.New()
	}

	return nil
}
