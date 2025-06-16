package controller

import (
	"net/http"
	"strconv"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
)

// ProductController struct will hold the repository instance
type ProductController struct {
	ProductRepo  repository.ProductRepository
	CategoryRepo repository.CategoryRepository
	PriceHistoryRepo repository.PriceHistoryRepository
}

// NewProductController creates a new ProductController instance
func NewProductController(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, priceHistoryRepo repository.PriceHistoryRepository) *ProductController {
	return &ProductController{
		ProductRepo:  productRepo,
		CategoryRepo: categoryRepo,
		PriceHistoryRepo: priceHistoryRepo,
	}
}

// GetProducts returns all products
func (pc *ProductController) GetProducts(c *gin.Context) {
	// Get all products from the repository
	prods, err := pc.ProductRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the products as JSON
	c.JSON(http.StatusOK, prods)
}

// CreateProduct adds a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req forms.ProductForm

	// Bind the incoming JSON to the product request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := pc.CategoryRepo.FindByID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If category is not found, return an error
	if category == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}

	// Create the product instance
	product := models.Product{
		Name:       req.Name,
		Stock:      req.Stock,
		Price:      req.Price,
		Location:   req.Location,
		CategoryID: req.CategoryID,
		Category:   *category, // Associate the category directly
	}

	// Create the product using the repository
	if err := pc.ProductRepo.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting row"})
		return
	}

	// Return the created product
	c.JSON(http.StatusCreated, product)

}

func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	// Convert id to uint
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Find the product by ID
	product, err := pc.ProductRepo.FindByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching product"})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Convert id to uint
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req forms.ProductForm

	// Bind the incoming JSON to the product request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := pc.CategoryRepo.FindByID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}

	// Check if requested price is equal to current price
	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}

	// Find the existing product to compare prices
	existingProduct, err := pc.ProductRepo.FindByID(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching existing product"})
		return
	}

	if existingProduct == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if requested price is equal to current price
	if req.Price != existingProduct.Price {
		// Create price history entry
		priceHistory := models.PriceHistory{
			ProductID: existingProduct.ID,
			OldPrice:  existingProduct.Price,
			NewPrice:  req.Price,
		}
		if err := pc.PriceHistoryRepo.Create(&priceHistory); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating price history"})
			return
		}
	}

	product := models.Product{
		ID:         uint(productID),
		Name:       req.Name,
		Stock:      req.Stock,
		Price:      req.Price,
		Location:   req.Location,
		CategoryID: req.CategoryID,
	}

	if err := pc.ProductRepo.Update(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Convert id to uint
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Delete the product by ID
	if err := pc.ProductRepo.Delete(uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}


func (pc *ProductController) SalesReport(c *gin.Context) {

	// Parse query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 0 {
		page = 1
	}
	
	if limit < 0 {
		limit = 10
	}

	// Get paginated sales report
	report, total, err := pc.ProductRepo.FindAllWithSales(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":       report,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	})

}