package repository

import (
	"fmt"

	"backend-task/internal/constants"
	"backend-task/internal/user/models"
	"backend-task/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Group Repository Interface :
type GroupRepository interface {
	FindAllocatableGroupTx(gormDB *gorm.DB, base string) (*models.Group, error)
	IncrementGroupCountTx(gormDB *gorm.DB, name string) error
}

// GroupRepositoryDB Implementation :
type GroupRepositoryDB struct {
	gormDB *gorm.DB
}

// Constructor :
func NewGroupRepository(db *gorm.DB) GroupRepository {

	return &GroupRepositoryDB{gormDB: db}
}

func (groupRepositoryDB *GroupRepositoryDB) FindAllocatableGroupTx(gormDB *gorm.DB, base string) (*models.Group, error) {

	var group models.Group

	// Try To find Existing Group With Available Capacity.
	err := gormDB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("base = ? AND member_count < ?", base, constants.GroupCapacity).
		Order("index ASC").
		First(&group).Error

	if err == nil {

		return &group, nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {

		return nil, fmt.Errorf("%w: %v", utils.ErrFailedToFindGroup, err)
	}

	// No Available Group, Create New.
	var maxIndex int
	if err2 := gormDB.Model(&models.Group{}).Where("base = ?", base).Select("COALESCE(MAX(\"index\"),0)").Scan(&maxIndex).Error; err2 != nil {

		return nil, fmt.Errorf("%w: %v", utils.ErrFailedToGetMaxGroupIdx, err2)
	}

	group = models.Group{

		Base:     base,
		Index:    maxIndex + 1,
		Capacity: constants.GroupCapacity,
		Name:     fmt.Sprintf("%s-%d", base, maxIndex+1),
	}

	if err3 := gormDB.Create(&group).Error; err3 != nil {

		return nil, fmt.Errorf("%w: %v", utils.ErrFailedToCreateNewGroup, err3)
	}

	return &group, nil
}

func (groupRepositoryDB *GroupRepositoryDB) IncrementGroupCountTx(tx *gorm.DB, name string) error {

	return tx.Model(&models.Group{}).
		Where("name = ? AND member_count < ?", name, constants.GroupCapacity).
		Update("member_count", gorm.Expr("member_count + 1")).Error
}
