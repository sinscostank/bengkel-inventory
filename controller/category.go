package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
)

// CategoryController struct will hold the repository instance
type CategoryController struct {
	CategoryRepo repository.CategoryRepository
}

// NewCategoryController creates a new CategoryController instance
func NewCategoryController(categoryRepo repository.CategoryRepository) *CategoryController {
	return &CategoryController{
		CategoryRepo: categoryRepo,
	}
}

// GetCategories returns all categories
func (pc *CategoryController) GetCategories(c *gin.Context) {
	// Get all categories from the repository
	cats, err := pc.CategoryRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the categories as JSON
	c.JSON(http.StatusOK, cats)
}

// CreateCategory adds a new product
func (pc *CategoryController) CreateCategory(c *gin.Context) {
	var req forms.CategoryForm


	// Bind the incoming JSON to the product request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the product instance
	category := models.Category{
		Name:       req.Name,
	}

	// Create the product using the repository
	if err := pc.CategoryRepo.Create(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting row"})
		return
	}

	// Return the created product
	c.JSON(http.StatusCreated, category)
}