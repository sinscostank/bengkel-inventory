package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
)

// ActivityItemRepository defines methods to interact with the products table.
type ActivityItemRepository interface {
	Create(product *models.ActivityItem) error
	CreateMultiple(ActivityItems []*models.ActivityItem) error
}

// ActivityItemRepositoryImpl is the implementation of the ActivityItemRepository interface.
type ActivityItemRepositoryImpl struct {
	DB *gorm.DB
}

// NewActivityItemRepository creates a new instance of ActivityItemRepositoryImpl
func NewActivityItemRepository(db *gorm.DB) ActivityItemRepository {
	return &ActivityItemRepositoryImpl{
		DB: db,
	}
}

func (r *ActivityItemRepositoryImpl) Create(ActivityItem *models.ActivityItem) error {
	return r.DB.Create(ActivityItem).Error
}

func (r *ActivityItemRepositoryImpl) CreateMultiple(ActivityItems []*models.ActivityItem) error {
	if len(ActivityItems) == 0 {
		return nil // No items to create
	}

	return r.DB.Create(ActivityItems).Error
}