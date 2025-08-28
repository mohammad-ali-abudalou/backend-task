package services

import (

    "backend-task/internal/models"
    "backend-task/internal/repository"
    "backend-task/pkg/utils"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "fmt"
	"time"
)

type UserService struct {

    Repo *repository.UserRepository
    DB   *gorm.DB
}

func NewUserService(repo *repository.UserRepository, db *gorm.DB) *UserService {

    return &UserService{Repo: repo, DB: db}
}

// Assign group based on age and capacity
func (s *UserService) AssignGroup(age int) (string, error) {

    var base string
    switch {
    case age <= 12:
        base = "child"
    case age <= 17:
        base = "teen"
    case age <= 64:
        base = "adult"
    default:
        base = "senior"
    }

    // Concurrency-safe group assignment
    var groupName string
    err := s.DB.Transaction(func(tx *gorm.DB) error {
	
        var count int64
        for i := 1; ; i++ {
		
            candidate := base
            if i > 1 {
			
                candidate = fmt.Sprintf("%s-%d", base, i)
            }
			
            tx.Model(&models.User{}).Where("'group' = ?", candidate).Count(&count)
            if count < 3 {
			
                groupName = candidate
                break
            }
        }
		
        return nil
    })
	
    return groupName, err
}

func (s *UserService) CreateUser(name, email string, dob string) (*models.User, error) {

	fmt.Println(name)
	fmt.Println(email)
	fmt.Println(dob)

    if !utils.ValidateEmail(email) {
	
        return nil, fmt.Errorf("invalid email")
    }

    dateOfBirth, err := time.Parse("2006-01-02", dob)
    if err != nil {
	
        return nil, fmt.Errorf("invalid date_of_birth format")
    }

    if err := utils.ValidateDOB(dateOfBirth); err != nil {
	
        return nil, err
    }

    age := utils.CalculateAge(dateOfBirth)
    group, err := s.AssignGroup(age)
    if err != nil {
	
        return nil, err
    }

    user := &models.User{
	
        ID:          uuid.New(),
        Name:        name,
        Email:       email,
        DateOfBirth: dateOfBirth,
        Group:       group,
    }

    if err := s.Repo.Create(user); err != nil {
	
        return nil, err
    }

    return user, nil
}

func (s *UserService) UpdateUser(user *models.User, name, email string) error {

    if email != "" {
	
        if !utils.ValidateEmail(email) {
		
            return fmt.Errorf("invalid email")
        }
        user.Email = email
    }
	
    if name != "" {
        user.Name = name
    }
	
    return s.Repo.Update(user)
}
