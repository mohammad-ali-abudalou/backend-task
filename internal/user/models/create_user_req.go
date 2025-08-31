package models

type CreateUserReq struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}
