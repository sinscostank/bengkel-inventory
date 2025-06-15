package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
)

// ProductRepository defines methods to interact with the products table.
type ProductRepository interface {
	Create(product *models.Product) error
	FindAll() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	FindByIDs(ids []uint) ([]models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
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

// Create adds a new product to the database.
func (r *ProductRepositoryImpl) Create(product *models.Product) error {
	return r.DB.Create(product).Error
}

// FindAll fetches all products from the database.
func (r *ProductRepositoryImpl) FindAll() ([]models.Product, error) {
	var products []models.Product
	if err := r.DB.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryImpl) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.DB.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) Update(product *models.Product) error {
	return r.DB.Save(product).Error
}

func (r *ProductRepositoryImpl) Delete(id uint) error {
	var product models.Product
	if err := r.DB.First(&product, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&product).Error
}

func (r *ProductRepositoryImpl) FindByIDs(ids []uint) ([]models.Product, error) {
	var products []models.Product
	if err := r.DB.Preload("Category").Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}