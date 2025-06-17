package controller

import (
	"net/http"
	"strconv"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/service"
)

// ProductController struct holds the service instance
type ProductController struct {
	ProductService service.ProductService
}

// NewProductController creates a new ProductController instance
func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{
		ProductService: productService,
	}
}

// GetProducts returns all products
func (pc *ProductController) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	prods, total, err := pc.ProductService.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":          prods,
		"current_page":  page,
		"limit":         limit,
		"total_items":   total,
		"total_pages":   totalPages,
	})
}

// CreateProduct adds a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req forms.ProductForm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.ProductService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProductByID returns a product by ID
func (pc *ProductController) GetProductByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	product, err := pc.ProductService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates an existing product
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req forms.ProductForm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.ProductService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := pc.ProductService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// SalesReport generates a sales report for products
func (pc *ProductController) SalesReport(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	report, total, err := pc.ProductService.GetSalesReport(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales report"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":         report,
		"current_page": page,
		"limit":        limit,
		"total_items":  total,
		"total_pages":  totalPages,
	})
}
