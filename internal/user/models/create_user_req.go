package models

// Create User Req Represents The Payload For Creating A User :
type CreateUserReq struct {
	Name        string `json:"name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john@example.com"`
	DateOfBirth string `json:"date_of_birth" binding:"required" example:"1990-01-01"`
}
