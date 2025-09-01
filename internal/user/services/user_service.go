package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend-task/internal/constants"
	"backend-task/internal/user/models"
	"backend-task/internal/user/repository"
	"backend-task/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(name, email, dob string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(id string, name, email *string) (*models.User, error)
	ListUsersByFilter(group string) ([]models.User, error)
}

type userService struct {
	db     *gorm.DB
	users  repository.UserRepository
	groups repository.GroupRepository
}

func NewUserService(db *gorm.DB, users repository.UserRepository, groups repository.GroupRepository) UserService {

	return &userService{db: db, users: users, groups: groups}
}

// ---------------- Create User ----------------

func (userService *userService) CreateUser(name, email, dob string) (*models.User, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {

		return nil, utils.NewBadRequest(utils.ErrNameIsRequired.Error())
	}

	if !utils.ValidateEmail(email) {

		return nil, utils.NewBadRequest(utils.ErrInvalidEmailFormat.Error())
	}

	birth, err := time.Parse("2006-01-02", dob)
	if err != nil {

		return nil, utils.NewBadRequest(utils.ErrDateOfBirthFormat.Error())
	}

	if err := utils.ValidateDateOfBirth(birth); err != nil {

		return nil, err
	}

	// Ensure Email Uniqueness :
	exists, err := userService.users.IsEmailExists(context.Background(), email)
	if err != nil {

		return nil, err
	}

	if exists {

		return nil, utils.NewBadRequest(utils.ErrEmailAlreadyExists.Error())
	}

	baseGroup := ageToBaseGroup(birth)
	var createdUser *models.User

	// Transaction For Safe Group Assignment :
	err = userService.db.Transaction(func(tx *gorm.DB) error {

		group, err := userService.groups.FindAllocatableGroupTx(tx, baseGroup)
		if err != nil {

			return err
		}

		user := &models.User{

			Name:        name,
			Email:       email,
			DateOfBirth: birth,
			Group:       group.Name,
		}

		if err := userService.users.CreateNewUser(context.Background(), user); err != nil {

			return err
		}

		if err := userService.groups.IncrementGroupCountTx(tx, group.Name); err != nil {

			return err
		}

		createdUser = user
		return nil
	})

	if err != nil {

		return nil, err
	}

	return createdUser, nil
}

// ---------------- Get User ----------------

func (userService *userService) GetUserByID(id string) (*models.User, error) {

	uid, err := uuid.Parse(id)
	if err != nil {

		return nil, utils.NewBadRequest(utils.ErrUserNotFound.Error())
	}

	user, err := userService.users.GetUserByID(context.Background(), uid)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, utils.NewBadRequest(utils.ErrUserNotFound.Error())
		}

		return nil, err
	}

	return user, nil
}

// ---------------- Update User ----------------

func (userService *userService) UpdateUser(id string, name, email *string) (*models.User, error) {

	uid, err := uuid.Parse(id)
	if err != nil {

		return nil, utils.NewBadRequest(utils.ErrUserNotFound.Error())
	}

	user, err := userService.users.GetUserByID(context.Background(), uid)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, utils.NewBadRequest(utils.ErrUserNotFound.Error())
		}

		return nil, err
	}

	changed := false
	if name != nil {

		newName := strings.TrimSpace(*name)
		if newName == "" {

			return nil, utils.NewBadRequest(utils.ErrNameCannotBeEmpty.Error())
		}

		user.Name = newName
		changed = true
	}

	if email != nil {

		newEmail := strings.ToLower(strings.TrimSpace(*email))
		if !utils.ValidateEmail(newEmail) {

			return nil, utils.NewBadRequest(utils.ErrInvalidEmailFormat.Error())
		}

		if newEmail != user.Email {

			exists, err := userService.users.IsEmailExists(context.Background(), newEmail)
			if err != nil {

				return nil, err
			}

			if exists {

				return nil, utils.NewBadRequest(utils.ErrEmailAlreadyExists.Error())
			}

			user.Email = newEmail
			changed = true
		}
	}

	if !changed {

		return user, nil
	}

	if err := userService.users.UpdateUser(context.Background(), user, "name", "email"); err != nil {

		return nil, err
	}

	return user, nil
}

// ---------------- List Users ----------------

func (userService *userService) ListUsersByFilter(group string) ([]models.User, error) {

	return userService.users.ListUsers(context.Background(), group)
}

// ---------------- Helper ----------------

func ageToBaseGroup(birth time.Time) string {

	age := utils.CalculateAge(birth)
	switch {
	case age >= 0 && age <= 12:
		return constants.BaseGroupChild

	case age >= 13 && age <= 17:
		return constants.BaseGroupTeen

	case age >= 18 && age <= 64:
		return constants.BaseGroupAdult

	case age >= 65:
		return constants.BaseGroupSenior

	default:
		return constants.BaseGroupUnset
	}
}
