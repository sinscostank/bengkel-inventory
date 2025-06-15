package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
)

// CategoryRepository defines methods to interact with the products table.
type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	Create(product *models.Category) error
	FindByID(id uint) (*models.Category, error)
	// You can add other methods like FindByID, Update, Delete if needed
}

// CategoryRepositoryImpl is the implementation of the CategoryRepository interface.
type CategoryRepositoryImpl struct {
	DB *gorm.DB
}

// NewCategoryRepository creates a new instance of CategoryRepositoryImpl
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{
		DB: db,
	}
}

// FindAll fetches all products from the database.
func (r *CategoryRepositoryImpl) FindAll() ([]models.Category, error) {
	var categories []models.Category
	if err := r.DB.Preload("Products").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID fetches a category by its ID from the database.
func (r *CategoryRepositoryImpl) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.DB.First(&category, id).Error; err != nil {
		// If the category is not found, return nil and the error
		if err == gorm.ErrRecordNotFound {
			return nil, nil // or return an appropriate error
		}
		return nil, err
	}
	return &category, nil
}

// Create adds a new product to the database.
func (r *CategoryRepositoryImpl) Create(product *models.Category) error {
	return r.DB.Create(product).Error
}
