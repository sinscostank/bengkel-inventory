package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
)

// PriceHistoryRepository defines methods to interact with the products table.
type PriceHistoryRepository interface {
	Create(product *models.PriceHistory) error
}

// PriceHistoryRepositoryImpl is the implementation of the PriceHistoryRepository interface.
type PriceHistoryRepositoryImpl struct {
	DB *gorm.DB
}

// NewPriceHistoryRepository creates a new instance of PriceHistoryRepositoryImpl
func NewPriceHistoryRepository(db *gorm.DB) PriceHistoryRepository {
	return &PriceHistoryRepositoryImpl{
		DB: db,
	}
}

// Create adds a new product to the database.
func (r *PriceHistoryRepositoryImpl) Create(product *models.PriceHistory) error {
	return r.DB.Create(product).Error
}
