package models

import "time"

// Group Is Aa Bookkeeping Table To Safely Distribute User Capacity :
//
// @Description Defines A User Group ( e.g., child-1, teen-2, adult-3, etc. )
// That Is Used To Manage Capacity Limits For Different Age Categories.
type Group struct {

	// Group Name ( e.g., "adult-1", "senior-2" ).
	// @Required
	Name string `gorm:"primaryKey;size:64" json:"name" example:"adult-1"`

	// Base Age Category ( child, teen, adult, senior ).
	// @Required
	Base string `gorm:"not null;index;size:32" json:"base" example:"adult"`

	// Sequential Index For This Base Category ( e.g., 1, 2, 3 ).
	// @Required
	Index int `gorm:"not null;index" json:"index" example:"1"`

	// Maximum Capacity Of Users Allowed In This Group.
	// @Required
	Capacity int `gorm:"not null;default:3" json:"capacity" example:"3"`

	// Current Number Of Users Assigned To This Group.
	// @Required
	MemberCount int `gorm:"not null;default:0" json:"member_count" example:"2"`

	// Timestamp When The Group Was Created.
	CreatedAt time.Time `json:"created_at" example:"2025-09-01T12:00:00Z"`

	// Timestamp When The Group Was Last Updated.
	UpdatedAt time.Time `json:"updated_at" example:"2025-09-01T12:30:00Z"`
}
