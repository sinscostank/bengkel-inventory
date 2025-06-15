package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
	"github.com/sinscostank/bengkel-inventory/utils"
)

// ActivityController struct will hold the repository instance
type ActivityController struct {
	ActivityRepo repository.ActivityRepository
	ActivityItemRepo repository.ActivityItemRepository
	ProductRepo repository.ProductRepository
	StockTransactionRepo repository.StockTransactionRepository
}

// NewActivityController creates a new ActivityController instance
func NewActivityController(ActivityRepo repository.ActivityRepository, ActivityItemRepo repository.ActivityItemRepository, StockTransactionRepo repository.StockTransactionRepository, ProductRepo repository.ProductRepository) *ActivityController {
	return &ActivityController{
		ActivityRepo: ActivityRepo,
		ActivityItemRepo: ActivityItemRepo,
		ProductRepo: ProductRepo,
		StockTransactionRepo: StockTransactionRepo,
	}
}

// GetCategories returns all categories
func (pc *ActivityController) GetActivities(c *gin.Context) {

	// Get all categories from the repository
	cats, err := pc.ActivityRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the categories as JSON
	c.JSON(http.StatusOK, cats)
}

// CreateActivity adds a new product
func (pc *ActivityController) CreateActivity(c *gin.Context) {
	
	var req forms.ActivityForm

	// Modify req.Type based on the route
	switch c.FullPath() {
	case "/activities":
		req.Type = "outbound" // If the request is to this route
	case "/stock-transactions":
		req.Type = "inbound" // If the request is to this route
	}

	claimsRaw, exists := c.Get("userClaims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

	userClaims, ok := claimsRaw.(*utils.UserClaims)  // or just UserClaims if you stored it by value
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user claims type"})
        return
    }

    userID := userClaims.ID
	userRole := userClaims.Role

	// Check activity type (inbound or outbound)
	if req.Type != "inbound" && req.Type != "outbound" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity type. Must be 'inbound' or 'outbound'."})
		return
	}

	if req.Type == "inbound" {
		// If the activity type is inbound, ensure the user is an admin
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create inbound activities"})
			return
		}
	}


	// Bind the incoming JSON to the product request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputMap := make(map[uint]uint) // key = product ID, value = requested quantity

	for _, item := range req.Products {
		inputMap[item.ID] = item.Quantity
	}

	productIDs := make([]uint, len(req.Products))
	for i, product := range req.Products {
		productIDs[i] = uint(product.ID)
	}

	// Validate that there are no duplicate product IDs
	productIDSet := make(map[uint]struct{})
	for _, productID := range productIDs {
		if _, exists := productIDSet[productID]; exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate product ID found: " + strconv.Itoa(int(productID))})
			return
		}
		productIDSet[productID] = struct{}{}
	}

	products, err := pc.ProductRepo.FindByIDs(productIDs)	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) != len(req.Products) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some products not found"})
		return
	}

	type ProductWithQuantity struct {
		ID       uint    `json:"id"`
		Name     string  `json:"name"`
		Stock    int     `json:"stock"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}

	var result []ProductWithQuantity

	for _, product := range products {
		requestedQty, found := inputMap[uint(product.ID)]
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID not found in input: " + strconv.Itoa(int(product.ID))})
			return
		}
	
		if product.Stock < int(requestedQty) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock for product ID: " + strconv.Itoa(int(product.ID)),
			})
			return
		}

		result = append(result, ProductWithQuantity{
			ID:       product.ID,
			Name:     product.Name,
			Stock:    product.Stock,
			Quantity: int(requestedQty),
			Price:    product.Price,
		})
	}

	// Create the product instance
	activity := models.Activity{
		UserID: userID,
		Status: "success", // Assuming status is always 'success' for creation
		Type:  req.Type,
	}

	// Create the activity using the repository
	if err := pc.ActivityRepo.Create(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting row"})
		return
	}

	//Create the activity items
	var activity_items []models.ActivityItem
	for _, product := range result {
		item := models.ActivityItem{
			ActivityID:     activity.ID,
			ProductID:      uint(product.ID),
			Quantity:       int(product.Quantity),
			PriceAtTime:    product.Price,
			DiscountAmount: 0,
			FinalPrice:     product.Price - 0, // Assuming no discount for simplicity
		}
		activity_items = append(activity_items, item)
	}

	// Create pointer for activity items
	activityItemPointers := make([]*models.ActivityItem, len(activity_items))
	for i := range activity_items {
		activityItemPointers[i] = &activity_items[i]
	}

	// Save the activity items
	if err := pc.ActivityItemRepo.CreateActivityItems(activityItemPointers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting activity items"})
		return
	}

	// Create stock transactions for each product
	var stock_transactions []models.StockTransaction
	
	for _, activity_item := range activityItemPointers {

		ChangeQuantity := 0

		if req.Type == "inbound" {
			ChangeQuantity = activity_item.Quantity // Positive for stock addition
		} else if req.Type == "outbound" {
			ChangeQuantity = -activity_item.Quantity // Negative for stock deduction
		}

		transaction := models.StockTransaction{
			ProductID: uint(activity_item.ProductID),
			ChangeQuantity:  ChangeQuantity, // Negative for stock deduction, positive for addition
			ActivityItemID: &activity_item.ID,
			Note: "Stock deduction for activity",
		}
		stock_transactions = append(stock_transactions, transaction)
	}

	// Create pointer for stock transactions
	stockTransactionPointers := make([]*models.StockTransaction, len(stock_transactions))
	for i := range stock_transactions {
		stockTransactionPointers[i] = &stock_transactions[i]
	}

	// Save the activity items
	if err := pc.StockTransactionRepo.CreateStockTransactions(stockTransactionPointers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting stock transactions"})
		return
	}

	// Update the products stock
	for _, product := range products {
		requestedQty, found := inputMap[uint(product.ID)]
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID not found in input: " + strconv.Itoa(int(product.ID))})
			return
		}
		
		if req.Type == "inbound" {
			product.Stock += int(requestedQty) // Increase stock for inbound
		} else if req.Type == "outbound" {
			product.Stock -= int(requestedQty)
		}

		// Update the product stock in the repository
		if err := pc.ProductRepo.Update(&product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating product stock"})
			return
		}
	}

	// Return the created activity with its items
	c.JSON(http.StatusCreated, activity)

}

func (pc *ActivityController) GetActivityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Activity ID"})
		return
	}

	Activity, err := pc.ActivityRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if Activity == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Activity not found"})
		return
	}

	c.JSON(http.StatusOK, Activity)
}