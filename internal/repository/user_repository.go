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
	Create(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, u *models.User, fields ...string) error
	List(ctx context.Context, group string) ([]models.User, error)
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

func (r *UserRepositoryDB) Create(ctx context.Context, u *models.User) error {

	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepositoryDB) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {

	var u models.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {

		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryDB) Update(ctx context.Context, u *models.User, fields ...string) error {

	return r.db.WithContext(ctx).Model(u).Select(fields).Updates(u).Error
}

func (r *UserRepositoryDB) List(ctx context.Context, group string) ([]models.User, error) {

	var users []models.User

	q := r.db.WithContext(ctx).Order("created_at asc")
	if group != "" {
		q = q.Where("\"group\" = ?", group)
	}

	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryDB) IsEmailExists(ctx context.Context, email string) (bool, error) {

	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Using Row-Level Locking to Implement Group Repositories.
func (r *GroupRepositoryDB) FindAllocatableGroupTx(tx *gorm.DB, base string) (*models.Group, error) {

	var g models.Group
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("base = ? AND member_count < 3 ", base).Order("index ASC").First(&g).Error
	if err == nil {
		return &g, nil
	}

	if err == gorm.ErrRecordNotFound {

		// Get Max Index.
		var maxIndex int
		if err2 := tx.Model(&models.Group{}).Where("base = ?", base).Select("COALESCE(MAX(\"index\"), 0)").Scan(&maxIndex).Error; err2 != nil {

			return nil, err2
		}

		g = models.Group{

			Base:     base,
			Index:    maxIndex + 1,
			Capacity: 3,
			Name:     fmtGroupName(base, maxIndex+1),
		}

		if err3 := tx.Create(&g).Error; err3 != nil {

			return nil, err3
		}

		return &g, nil
	}

	return nil, err
}

func (r *GroupRepositoryDB) IncrementGroupCountTx(tx *gorm.DB, name string) error {

	return tx.Model(&models.Group{}).
		Where("name = ? AND member_count < 3 ", name).
		Update("member_count", gorm.Expr("member_count + 1")).Error
}

func fmtGroupName(base string, idx int) string {

	return base + "-" + fmtInt(idx)
}

func fmtInt(n int) string {
	return fmt.Sprintf("%d", n)
}
