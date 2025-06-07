package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/config"
	"github.com/sinscostank/bengkel-inventory/entity"
)

// GetCategories returns all categories
func GetCategories(c *gin.Context) {
	var cats []entity.Category
	if err := config.DB.Find(&cats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	var cat entity.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}
