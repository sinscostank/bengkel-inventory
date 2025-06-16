package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/dto"
	"gorm.io/gorm"
	"fmt"
)

// ProductRepository defines methods to interact with the products table.
type ProductRepository interface {
	Create(product *models.Product) error
	FindAll() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	FindByIDs(ids []uint) ([]models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	FindAllWithSales(page int, limit int) ([]dto.ProductSalesDTO, int64, error)
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

func (r *ProductRepositoryImpl) FindAllWithSales(page int, limit int) ([]dto.ProductSalesDTO, int64, error) {
	var result []dto.ProductSalesDTO
	var total int64

	offset := (page - 1) * limit

	// Count total products
	err := r.DB.Model(&models.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			p.id, 
			p.name, 
			p.stock,
			c.category_name,
			ABS(COALESCE(SUM(ai.quantity), 0)) as total_sales
		FROM products p
		JOIN categories c ON p.category_id = c.id
		LEFT JOIN activity_items ai ON ai.product_id = p.id AND ai.quantity < 0
		WHERE p.deleted_at IS NULL
		GROUP BY p.id, p.name, p.stock
		ORDER BY total_sales DESC
		LIMIT ? OFFSET ?
	`

	
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	if err := r.DB.Raw(query, limit, offset).Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}