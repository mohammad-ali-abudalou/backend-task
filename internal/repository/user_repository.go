package repository

import (
	"context"
	"fmt"

	"backend-task/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, u *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, u *models.User, fields ...string) error
	ListUsers(ctx context.Context, group string) ([]models.User, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
}

type GroupRepository interface {
	FindAllocatableGroupTx(tx *gorm.DB, base string) (*models.Group, error)
	IncrementGroupCountTx(tx *gorm.DB, name string) error
}

type UserRepositoryDB struct {
	db *gorm.DB
}

type GroupRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {

	return &UserRepositoryDB{db: db}
}

func NewGroupRepository(db *gorm.DB) GroupRepository {

	return &GroupRepositoryDB{db: db}
}

func (r *UserRepositoryDB) CreateNewUser(context context.Context, user *models.User) error {

	return r.db.WithContext(context).Create(user).Error
}

func (r *UserRepositoryDB) GetUserByID(context context.Context, userId uuid.UUID) (*models.User, error) {

	var u models.User
	if err := r.db.WithContext(context).First(&u, "id = ?", userId).Error; err != nil {

		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryDB) UpdateUser(context context.Context, user *models.User, fields ...string) error {

	return r.db.WithContext(context).Model(user).Select(fields).Updates(user).Error
}

func (r *UserRepositoryDB) ListUsers(context context.Context, group string) ([]models.User, error) {

	var users []models.User

	q := r.db.WithContext(context).Order("created_at asc")
	if group != "" {
		q = q.Where("\"group\" = ?", group)
	}

	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryDB) IsEmailExists(context context.Context, email string) (bool, error) {

	var count int64
	if err := r.db.WithContext(context).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Using Row-Level Locking to Implement Group Repositories.
func (groupRepository *GroupRepositoryDB) FindAllocatableGroupTx(db *gorm.DB, baseGroup string) (*models.Group, error) {

	var group models.Group
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("base = ? AND member_count < 3 ", baseGroup).Order("index ASC").First(&group).Error
	if err == nil {
		return &group, nil
	}

	if err == gorm.ErrRecordNotFound {

		// Get Max Index.
		var maxIndex int
		if err2 := db.Model(&models.Group{}).Where("base = ?", baseGroup).Select("COALESCE(MAX(\"index\"), 0)").Scan(&maxIndex).Error; err2 != nil {

			return nil, err2
		}

		group = models.Group{

			Base:     baseGroup,
			Index:    maxIndex + 1,
			Capacity: 3,
			Name:     fmtGroupName(baseGroup, maxIndex+1),
		}

		if err3 := db.Create(&group).Error; err3 != nil {

			return nil, err3
		}

		return &group, nil
	}

	return nil, err
}

func (groupRepository *GroupRepositoryDB) IncrementGroupCountTx(db *gorm.DB, name string) error {

	return db.Model(&models.Group{}).
		Where("name = ? AND member_count < 3 ", name).
		Update("member_count", gorm.Expr("member_count + 1")).Error
}

func fmtGroupName(base string, idx int) string {

	return base + "-" + fmtInt(idx)
}

func fmtInt(n int) string {
	return fmt.Sprintf("%d", n)
}
