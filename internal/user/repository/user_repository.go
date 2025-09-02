package repository

import (
	"context"
	"fmt"

	"backend-task/internal/user/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Repository Interface :
type UserRepository interface {
	CreateNewUser(context context.Context, user *models.User) error
	GetUserByID(context context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(context context.Context, user *models.User, fields ...string) error
	ListUsers(context context.Context, group string) ([]*models.User, error)
	IsEmailExists(context context.Context, email string) (bool, error)
}

// UserRepositoryDB Implementation :
type UserRepositoryDB struct {
	gormDB *gorm.DB
}

// Constructor
func NewUserRepository(db *gorm.DB) UserRepository {

	return &UserRepositoryDB{gormDB: db}
}

func (userRepositoryDB *UserRepositoryDB) CreateNewUser(context context.Context, user *models.User) error {

	return userRepositoryDB.gormDB.WithContext(context).Create(user).Error
}

func (userRepositoryDB *UserRepositoryDB) GetUserByID(context context.Context, userID uuid.UUID) (*models.User, error) {

	var user models.User
	if err := userRepositoryDB.gormDB.WithContext(context).First(&user, "id = ?", userID).Error; err != nil {

		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (userRepositoryDB *UserRepositoryDB) UpdateUser(context context.Context, user *models.User, fields ...string) error {

	return userRepositoryDB.gormDB.WithContext(context).Model(user).Select(fields).Updates(user).Error
}

func (userRepositoryDB *UserRepositoryDB) ListUsers(context context.Context, group string) ([]*models.User, error) {

	var users []*models.User
	gormDB := userRepositoryDB.gormDB.WithContext(context).Order("created_at asc")
	if group != "" {

		gormDB = gormDB.Where("\"group\" = ?", group)
	}

	if err := gormDB.Find(&users).Error; err != nil {

		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (userRepositoryDB *UserRepositoryDB) IsEmailExists(context context.Context, email string) (bool, error) {

	var count int64
	if err := userRepositoryDB.gormDB.WithContext(context).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {

		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}
