package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
)

type ProductService interface {
	GetAll(page, limit int) ([]models.Product, int64, error)
	Create(req forms.ProductForm) (models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(id string, form forms.UpdateProductForm) (models.Product, error)
	Delete(id string) error
	GetSalesReport(page, limit int) ([]models.ProductSales, int64, error)
}

// ProductService struct holds the repository instance
type productService struct {
	ProductRepo repository.ProductRepository
	CategoryRepo repository.CategoryRepository
	PriceHistoryRepo repository.PriceHistoryRepository
}

// NewProductService creates a new ProductService instance
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, priceHistoryRepo repository.PriceHistoryRepository) ProductService {
	return &productService{
		ProductRepo:      productRepo,
		CategoryRepo:     categoryRepo,
		PriceHistoryRepo: priceHistoryRepo,
	}
}

// GetAll retrieves all products with pagination
func (ps *productService) GetAll(page, limit int) ([]models.Product, int64, error) {
	return ps.ProductRepo.FindAll(page, limit)
}

// Create adds a new product
func (ps *productService) Create(req forms.ProductForm) (models.Product, error) {
	category, err := ps.CategoryRepo.FindByID(req.CategoryID)
	if err != nil || category == nil {
		return models.Product{}, errors.New("invalid category ID")
	}

	product := models.Product{
		Name:       req.Name,
		Stock:      req.Stock,
		Price:      req.Price,
		Location:   req.Location,
		CategoryID: req.CategoryID,
		Category:   *category,
	}

	if err := ps.ProductRepo.Create(&product); err != nil {
		return models.Product{}, err
	}

	return product, nil
}

// GetByID retrieves a product by ID
func (ps *productService) GetByID(id uint) (*models.Product, error) {
	product, err := ps.ProductRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// Update modifies an existing product
func (ps *productService) Update(id string, req forms.UpdateProductForm) (models.Product, error) {
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return models.Product{}, errors.New("invalid product ID")
	}

	existingProduct, err := ps.ProductRepo.FindByID(uint(productID))
    if err != nil || existingProduct == nil {
        return models.Product{}, errors.New("product not found")
    }

	category, err := ps.CategoryRepo.FindByID(req.CategoryID)
	if err != nil || category == nil {
		return models.Product{}, errors.New("invalid category ID")
	}

	product := models.Product{
		ID:         uint(productID),
		Name:       req.Name,
		Stock:      existingProduct.Stock,
		Price:      req.Price,
		Location:   req.Location,
		CategoryID: req.CategoryID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Only log price history if price actually changed
    if existingProduct.Price != req.Price {
        history := models.PriceHistory{
            ProductID:   product.ID,
            OldPrice:    existingProduct.Price,
            NewPrice:    req.Price,
            DateChanged: time.Now(),
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        if err := ps.PriceHistoryRepo.Create(&history); err != nil {
            return models.Product{}, err
        }
    }

	if err := ps.ProductRepo.Update(&product); err != nil {
		return models.Product{}, err
	}

	return product, nil
}

// Delete removes a product by ID
func (ps *productService) Delete(id string) error {
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("invalid product ID")
	}
	return ps.ProductRepo.Delete(uint(productID))
}

// SalesReport returns paginated sales data per product
func (ps *productService) GetSalesReport(page, limit int) ([]models.ProductSales, int64, error) {
	return ps.ProductRepo.FindAllWithSales(page, limit)
}