package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
	"errors"
)

// CategoryRepository defines methods to interact with the products table.
type CategoryRepository interface {
	Create(product *models.Category) error
	FindAll(page int, limit int) ([]models.Category, int64, error)
	FindByID(id uint) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
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
func (r *CategoryRepositoryImpl) FindAll(page int, limit int) ([]models.Category, int64, error) {
	
	var categories []models.Category

	var total int64

	// Count total products
	err := r.DB.Model(&models.Category{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		if err := r.DB.Preload("Products").
			Limit(limit).
			Offset(offset).
			Find(&categories).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.DB.Preload("Products").Find(&categories).Error; err != nil {
			return nil, 0, err
		}
	}

	return categories, total, nil
}

// FindByID fetches a category by its ID from the database.
func (r *CategoryRepositoryImpl) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.DB.First(&category, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return nil, nil to indicate not found without error
		return nil, nil
	}

	if err != nil {
		// Real DB error
		return nil, err
	}
	
	return &category, nil
}

// Create adds a new product to the database.
func (r *CategoryRepositoryImpl) Create(product *models.Category) error {
	return r.DB.Create(product).Error
}

// Update updates an existing category in the database.
func (r *CategoryRepositoryImpl) Update(category *models.Category) error {
	return r.DB.Save(category).Error
}

// Delete removes a category from the database by its ID.
func (r *CategoryRepositoryImpl) Delete(id uint) error {
	var category models.Category
	if err := r.DB.First(&category, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&category).Error
}

