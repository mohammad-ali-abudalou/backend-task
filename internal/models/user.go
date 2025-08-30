package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Email       string    `json:"email" gorm:"not null;uniqueIndex;size:320"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"type:date;not null"`
	Group       string    `json:"group" gorm:"not null;index;size:64"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.ID == uuid.Nil {

		u.ID = uuid.New()
	}

	return nil
}

// Group Is A Bookkeeping Table That Uses Row-Level Locks To Safely Distribute Capacity.
type Group struct {
	Name        string `gorm:"primaryKey;size:64"`
	Base        string `gorm:"not null;index;size:32"` // child | teen | adult | senior
	Index       int    `gorm:"not null;index"`         // 1,2,3
	Capacity    int    `gorm:"not null;default:3"`
	MemberCount int    `gorm:"not null;default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
