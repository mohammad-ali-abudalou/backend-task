package models

import "time"

// Group Is Aa Bookkeeping Table Tto Safely Distribute User Capacity :
type Group struct {
	Name        string    `gorm:"primaryKey;size:64" json:"name"`
	Base        string    `gorm:"not null;index;size:32" json:"base"` // child | teen | adult | senior
	Index       int       `gorm:"not null;index" json:"index"`        // 1,2,3
	Capacity    int       `gorm:"not null;default:3" json:"capacity"` // max users per group
	MemberCount int       `gorm:"not null;default:0" json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
