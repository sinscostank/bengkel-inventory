package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/controller"
)

func SetupRoutes(
	userController *controller.UserController, 
	productController *controller.ProductController, 
	categoryController *controller.CategoryController,
	) *gin.Engine {

	r := gin.Default()

	// Health‐check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Category
	r.GET("/categories", categoryController.GetCategories)
	r.POST("/categories", categoryController.CreateCategory)

	// Product
	r.GET("/products", productController.GetProducts)
	r.POST("/products", productController.CreateProduct)

	// User
	r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)

	return r

}
