package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
)

// ProductRepository defines methods to interact with the products table.
type ProductRepository interface {
	FindAll() ([]models.Product, error)
	Create(product *models.Product) error
	// You can add other methods like FindByID, Update, Delete if needed
}

// ProductRepositoryImpl is the implementation of the ProductRepository interface.
type ProductRepositoryImpl struct {
	DB *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepositoryImpl
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

// FindAll fetches all products from the database.
func (r *ProductRepositoryImpl) FindAll() ([]models.Product, error) {
	var products []models.Product
	if err := r.DB.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Create adds a new product to the database.
func (r *ProductRepositoryImpl) Create(product *models.Product) error {
	return r.DB.Create(product).Error
}
