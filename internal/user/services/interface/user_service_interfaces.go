package serviceInterface

import "backend-task/internal/user/models"

// User Service Defines All Operations The Service Must Provide :
type UserService interface {

	// CreateUser Creates A User And Assigns Them To A Group Automatically.
	CreateUser(name, email, dob string) (*models.User, error)

	// GetUserByID Retrieves A User By UUID.
	GetUserByID(id string) (*models.User, error)

	// UpdateUser Updates The Name And/Or Email Of A User.
	UpdateUser(id string, name, email *string) (*models.User, error)

	// ListUsersByFilter Lists Users Optionally Filtered By Group.
	ListUsersByFilter(group string) ([]*models.User, error)
}
