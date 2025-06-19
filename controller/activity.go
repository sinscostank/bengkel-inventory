package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/utils"
	"github.com/sinscostank/bengkel-inventory/service"

)

// ActivityController struct will hold the repository instance
type ActivityController struct {
	ActivityService service.ActivityService
}

// NewActivityController creates a new ActivityController instance
func NewActivityController(ActivityService service.ActivityService) *ActivityController {
	return &ActivityController{
		ActivityService: ActivityService,
	}
}

// GetCategories returns all categories
func (pc *ActivityController) GetActivities(c *gin.Context) {

	// Get all categories from the repository
	acts, err := pc.ActivityService.GetAll()
	if err != nil {
		if err.Error() == "activity not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the categories as JSON
	c.JSON(http.StatusOK, gin.H{
		"data":         acts,
	})
}

// CreateActivity adds a new product
func (pc *ActivityController) CreateActivity(c *gin.Context) {
	var req forms.ActivityForm

	// Determine type based on route
	switch c.FullPath() {
	case "/activities":
		req.Type = "outbound"
	case "/stock-transactions":
		req.Type = "inbound"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route"})
		return
	}

	claimsRaw, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userClaims := claimsRaw.(*utils.UserClaims)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity, err := pc.ActivityService.Create(userClaims.ID, userClaims.Role, &req)
	if err != nil {
		if err.Error() == "invalid activity type" || err.Error() == "duplicate product ID found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "only admins can create inbound activities" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func (pc *ActivityController) GetActivityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Activity ID"})
		return
	}

	activity, err := pc.ActivityService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "activity not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, activity)
}