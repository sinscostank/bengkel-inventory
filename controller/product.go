package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
)

// ProductController struct will hold the repository instance
type ProductController struct {
	ProductRepo repository.ProductRepository
	CategoryRepo  repository.CategoryRepository
}

// NewProductController creates a new ProductController instance
func NewProductController(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) *ProductController {
	return &ProductController{
		ProductRepo:  productRepo,
		CategoryRepo: categoryRepo,
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
	if err!= nil {
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