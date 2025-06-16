package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/controller"
	"github.com/sinscostank/bengkel-inventory/middleware"
	"github.com/sinscostank/bengkel-inventory/repository"
	"gorm.io/gorm"

)

func SetupRoutes(
	dbConn *gorm.DB,
) *gin.Engine {

	// Create repository
	userRepo := repository.NewUserRepository(dbConn)
	productRepo := repository.NewProductRepository(dbConn)
	categoryRepo := repository.NewCategoryRepository(dbConn)
	activityRepo := repository.NewActivityRepository(dbConn)
	actiityItemRepo := repository.NewActivityItemRepository(dbConn)
	stockTransactionRepo := repository.NewStockTransactionRepository(dbConn)
	priceHistoryRepo := repository.NewPriceHistoryRepository(dbConn)

	// Create controllers
	userController := controller.NewUserController(userRepo)
	productController := controller.NewProductController(productRepo, categoryRepo, priceHistoryRepo)
	categoryController := controller.NewCategoryController(categoryRepo)
	activityController := controller.NewActivityController(activityRepo, actiityItemRepo, stockTransactionRepo, productRepo)


	// Initialize Gin router
	r := gin.Default()

	authenticatedGroup := r.Group("", middleware.AuthMiddleware())
	{
		// Category
		categoryGroup := authenticatedGroup.Group("/categories")
		{
			categoryGroup.GET("", categoryController.GetCategories)
			categoryGroup.GET("/:id", categoryController.GetCategoryByID)
	
			// Admin routes for categories
			adminCategoryGroup := categoryGroup.Group("", middleware.AdminMiddleware())
			{
				// Admin only routes
				adminCategoryGroup.POST("", categoryController.CreateCategory)
				adminCategoryGroup.PUT("/:id", categoryController.UpdateCategory)
				adminCategoryGroup.DELETE("/:id", categoryController.DeleteCategory)
			}
		
		}
	
		// Product
		productGroup := authenticatedGroup.Group("/products")
		{
			productGroup.GET("", productController.GetProducts)
			productGroup.GET("/:id", productController.GetProductByID)
		
			// Admin routes for products
			adminProductGroup := productGroup.Group("", middleware.AdminMiddleware())
			{
				adminProductGroup.POST("", productController.CreateProduct)
				adminProductGroup.PUT("/:id", productController.UpdateProduct)
				adminProductGroup.DELETE("/:id", productController.DeleteProduct)
			}
		}
	
		// Product
		activitiesGroup := authenticatedGroup.Group("/activities")
		{
			activitiesGroup.GET("", activityController.GetActivities)
			activitiesGroup.POST("", activityController.CreateActivity)
		}
	
		// Stock Transactions
		r.POST("/stock-transactions",  middleware.AdminMiddleware(), activityController.CreateActivity) 
	
		// Sales Report
		r.GET("/sales-report", productController.SalesReport)
	}

	// Health‐check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// User
	r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)


	return r

}
