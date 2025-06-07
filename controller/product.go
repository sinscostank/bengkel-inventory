package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/config"
	"github.com/sinscostank/bengkel-inventory/entity"
)

// GetProducts returns all products
func GetProducts(c *gin.Context) {
	var prods []entity.Product
	if err := config.DB.Preload("Category").Find(&prods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prods)
}

// CreateProduct adds a new product
func CreateProduct(c *gin.Context) {
	var prod entity.Product
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&prod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, prod)
}
