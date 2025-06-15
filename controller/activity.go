package controller

import (
	"fmt"
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

	fmt.Println("CreateActivity called")
	
	var req forms.ActivityForm

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

	// Bind the incoming JSON to the product request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputMap := make(map[int]int) // key = product ID, value = requested quantity

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
		requestedQty, found := inputMap[int(product.ID)]
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID not found in input: " + strconv.Itoa(int(product.ID))})
			return
		}
	
		if product.Stock < requestedQty {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock for product ID: " + strconv.Itoa(int(product.ID)),
			})
			return
		}

		result = append(result, ProductWithQuantity{
			ID:       product.ID,
			Name:     product.Name,
			Stock:    product.Stock,
			Quantity: requestedQty,
			Price:    product.Price,
		})
	}

	// Create the product instance
	activity := models.Activity{
		UserID:userID,
		Status: "success", // Assuming status is always 'success' for creation
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
		transaction := models.StockTransaction{
			ProductID: uint(activity_item.ProductID),
			ChangeQuantity:  -activity_item.Quantity, // Negative for stock deduction
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
		requestedQty, found := inputMap[int(product.ID)]
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID not found in input: " + strconv.Itoa(int(product.ID))})
			return
		}
		product.Stock -= requestedQty
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