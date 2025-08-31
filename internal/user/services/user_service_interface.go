package service

import (
	"backend-task/internal/user/models"
	"context"
)

type IUserService interface {
	CreateUser(ctx context.Context, u *models.User) error
	ListUsersByFilter(ctx context.Context, filter string) ([]models.User, error)
}
