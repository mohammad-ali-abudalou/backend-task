package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend-task/internal/user/models"
	repository "backend-task/internal/user/repository"
	"backend-task/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(name, email, dob string) (*models.User, error)
	GetUserById(id string) (*models.User, error)
	UpdateUser(id string, name, email *string) (*models.User, error)
	ListUsersByFilter(group string) ([]models.User, error)
}

type UserServiceStruct struct {
	db     *gorm.DB
	users  repository.UserRepository
	groups repository.GroupRepository
}

func NewUserService(db *gorm.DB, user repository.UserRepository, group repository.GroupRepository) UserService {

	return &UserServiceStruct{db: db, users: user, groups: group}
}

func (userService *UserServiceStruct) CreateUser(name, email, dob string) (*models.User, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {

		return nil, utils.NewBadRequest(utils.ErrNameIsRequired.Error())
	}

	if !utils.ValidateEmail(email) {

		return nil, utils.NewBadRequest(utils.ErrInvalidEmailFormat.Error())
	}

	// Parse DOB (YYYY-MM-DD).
	birth, err := time.Parse("2006-01-02", dob)
	if err != nil {

		return nil, utils.NewBadRequest(utils.ErrDateOfBirthFormat.Error())
	}

	if err := utils.ValidateDateOfBirth(birth); err != nil { // birthdate Unable To Occur In The Future.

		return nil, err
	}

	// Check If Email Is Exists.
	exists, err := userService.users.IsEmailExists(context.Background(), email)
	if err != nil {

		return nil, err
	}

	if exists {

		return nil, utils.NewBadRequest(utils.ErrEmailAlreadyExists.Error())
	}

	baseGroup := ageToBaseGroup(birth)

	var userCreated *models.User
	err = userService.db.Transaction(func(db *gorm.DB) error {

		group, err := userService.groups.FindAllocatableGroupTx(db, baseGroup)
		if err != nil {

			return err
		}

		user := &models.User{Name: name, Email: email, DateOfBirth: birth, Group: group.Name}
		if err := db.Create(user).Error; err != nil {

			return err
		}

		if err := userService.groups.IncrementGroupCountTx(db, group.Name); err != nil {

			return err
		}

		userCreated = user

		return nil
	})

	if err != nil {

		return nil, err
	}

	return userCreated, nil
}

func (userService *UserServiceStruct) GetUserById(id string) (*models.User, error) {

	userId, err := uuid.Parse(id)
	if err != nil {

		return nil, utils.NewBadRequest(utils.ErrUserNotFound.Error())
	}

	user, err := userService.users.GetUserByID(context.Background(), userId)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewBadRequest(utils.ErrUserNotFound.Error())
		}

		return nil, err
	}

	return user, nil
}

func (userService *UserServiceStruct) UpdateUser(id string, name, email *string) (*models.User, error) {

	userId, err := uuid.Parse(id)
	if err != nil {

		return nil, utils.NewNotFound(utils.ErrUserNotFound.Error())
	}

	user, err := userService.users.GetUserByID(context.Background(), userId)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, utils.NewNotFound(utils.ErrUserNotFound.Error())
		}

		return nil, err
	}

	changed := false
	if name != nil {

		newName := strings.TrimSpace(*name)
		if newName == "" {

			return nil, utils.NewBadRequest(utils.ErrNameCanNotEmpty.Error())
		}

		user.Name = newName
		changed = true
	}

	if email != nil {

		newEmail := strings.ToLower(strings.TrimSpace(*email))

		if !utils.ValidateEmail(newEmail) {
			return nil, utils.NewBadRequest(utils.ErrInvalidEmailFormat.Error())
		}

		// Email Changed, Ensure Uniqueness Email
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

func (userService *UserServiceStruct) ListUsersByFilter(group string) ([]models.User, error) {

	return userService.users.ListUsers(context.Background(), group)
}

func ageToBaseGroup(birth time.Time) string {

	age := utils.CalculateAge(birth)

	switch {

	case (age >= 0 && age <= 12):
		return utils.BaseGroupChild

	case (age >= 12 && age <= 17):
		return utils.BaseGroupTeen

	case (age >= 17 && age <= 64):
		return utils.BaseGroupAdult

	case age > 64:
		return utils.BaseGroupSenior

	default:
		return utils.BaseGroupUnset
	}
}
