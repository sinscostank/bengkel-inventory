package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/controller"
)

func SetupRoutes(userController *controller.UserController) *gin.Engine {

	r := gin.Default()

	// Health‐check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Category
	r.GET("/categories", controller.GetCategories)
	r.POST("/categories", controller.CreateCategory)

	// Product
	r.GET("/products", controller.GetProducts)
	r.POST("/products", controller.CreateProduct)

	// User
	r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)

	return r

}
