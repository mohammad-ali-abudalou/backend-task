package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"backend-task/internal/models"
	"backend-task/internal/repository"
	"backend-task/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(name, email, dob string) (*models.User, error)
	GetUser(id string) (*models.User, error)
	UpdateUser(id string, name, email *string) (*models.User, error)
	ListUsers(group string) ([]models.User, error)
}

type userService struct {
	db     *gorm.DB
	users  repository.UserRepository
	groups repository.GroupRepository
}

func NewUserService(db *gorm.DB, users repository.UserRepository, groups repository.GroupRepository) UserService {

	return &userService{db: db, users: users, groups: groups}
}

var emailRx = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

func (s *userService) CreateUser(name, email, dob string) (*models.User, error) {

	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {
		return nil, newBadRequest("Name Is Required !")
	}

	if !utils.ValidateEmail(email) {
		return nil, newBadRequest("Invalid Email Format !")
	}

	// Parse DOB (YYYY-MM-DD).
	birth, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return nil, newBadRequest("date_of_birth Must Be YYYY-MM-DD")
	}

	if err := utils.ValidateDOB(birth); err != nil { // date_of_birth Cannot Be In The Future.
		return nil, err
	}

	// Check Duplicate Email.
	exists, err := s.users.EmailExists(context.Background(), email)
	if err != nil {
		return nil, err
	}

	if exists {

		return nil, newBadRequest("Email Already Exists !")
	}

	base := ageToBaseGroup(birth)

	var created *models.User
	err = s.db.Transaction(func(tx *gorm.DB) error {

		// Lock / Select An Allocatable Group Row Or Create The Next One.
		g, err := s.groups.FindAllocatableGroupTx(tx, base)
		if err != nil {
			return err
		}

		u := &models.User{Name: name, Email: email, DateOfBirth: birth, Group: g.Name}
		if err := tx.Create(u).Error; err != nil {
			return err
		}

		if err := s.groups.IncrementGroupCountTx(tx, g.Name); err != nil {
			return err
		}

		created = u

		return nil
	})

	if err != nil {

		return nil, err
	}

	return created, nil
}

func (s *userService) GetUser(id string) (*models.User, error) {

	uuidV, err := uuid.Parse(id)
	if err != nil {
		return nil, newNotFound("User Not Found !")
	}

	u, err := s.users.GetByID(context.Background(), uuidV)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newNotFound("User Not Found !")
		}

		return nil, err
	}

	return u, nil
}

func (s *userService) UpdateUser(id string, name, email *string) (*models.User, error) {

	uuidV, err := uuid.Parse(id)
	if err != nil {

		return nil, newNotFound("User Not Found !")
	}

	u, err := s.users.GetByID(context.Background(), uuidV)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, newNotFound("User Not Found !")
		}

		return nil, err
	}

	changed := false
	if name != nil {

		n := strings.TrimSpace(*name)
		if n == "" {

			return nil, newBadRequest("Name Cannot Be Empty !")
		}

		u.Name = n
		changed = true
	}

	if email != nil {

		e := strings.ToLower(strings.TrimSpace(*email))
		if !emailRx.MatchString(e) {
			return nil, newBadRequest("Invalid Email Format !")
		}

		// If Email Changed, Ensure Uniqueness.
		if e != u.Email {

			exists, err := s.users.EmailExists(context.Background(), e)
			if err != nil {
				return nil, err
			}

			if exists {
				return nil, newBadRequest("Email Already Exists !")
			}

			u.Email = e
			changed = true
		}
	}

	if !changed {

		return u, nil
	}

	if err := s.users.Update(context.Background(), u, "name", "email"); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) ListUsers(group string) ([]models.User, error) {

	return s.users.List(context.Background(), group)
}

func ageToBaseGroup(birth time.Time) string {

	age := utils.CalculateAge(birth)
	switch {
	case (age >= 0 && age <= 12):
		return "child"
	case (age >= 12 && age <= 17):
		return "teen"
	case (age >= 17 && age <= 64):
		return "adult"
	case age > 64:
		return "senior"
	default:
		return "unset"
	}
}

// Errors :
type apiErr struct {
	Code int
	Msg  string
}

func (e apiErr) Error() string {

	return e.Msg
}

func newBadRequest(msg string) error {

	return apiErr{Code: 400, Msg: msg}
}

func newNotFound(msg string) error {

	return apiErr{Code: 404, Msg: msg}
}
