package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/controller"
	"github.com/sinscostank/bengkel-inventory/middleware"
)

func SetupRoutes(
	userController *controller.UserController,
	productController *controller.ProductController,
	categoryController *controller.CategoryController,
	activityController *controller.ActivityController,
) *gin.Engine {

	r := gin.Default()

	// Health‐check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Category
	r.GET("/categories", middleware.AuthMiddleware(), categoryController.GetCategories)
	r.POST("/categories", middleware.AuthMiddleware(), middleware.AdminMiddleware(), categoryController.CreateCategory)
	r.GET("/categories/:id", middleware.AuthMiddleware(), categoryController.GetCategoryByID)
	r.PUT("/categories/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), categoryController.UpdateCategory)
	r.DELETE("/categories/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), categoryController.DeleteCategory)

	// Product
	r.GET("/products", middleware.AuthMiddleware(), productController.GetProducts)
	r.POST("/products", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productController.CreateProduct)
	r.GET("/products/:id", middleware.AuthMiddleware(), productController.GetProductByID)
	r.PUT("/products/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productController.UpdateProduct)
	r.DELETE("/products/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productController.DeleteProduct)

	// User
	r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)

	// Activity
	r.GET("/activities", middleware.AuthMiddleware(), activityController.GetActivities)
	r.POST("/activities", middleware.AuthMiddleware(), activityController.CreateActivity)

	// Stock Transactions
	r.POST("/stock-transactions", middleware.AuthMiddleware(), middleware.AdminMiddleware(), activityController.CreateActivity) 

	return r

}
