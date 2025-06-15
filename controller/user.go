package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
	"github.com/sinscostank/bengkel-inventory/utils"
)

var userModel = new(models.User)
var userForm = new(forms.UserForm)

// UserController contains the repository for database access
type UserController struct {
	UserRepo *repository.UserRepository
}

// NewUserController creates a new instance of UserController
func NewUserController(repo *repository.UserRepository) *UserController {
	return &UserController{UserRepo: repo}
}

// RegisterUser handles user registration
func (uc *UserController) RegisterUser(c *gin.Context) {
	var req forms.RegisterForm

	// Parse the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		// This part handles the error if the binding fails
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"details": err.Error(), // Include more specific error details, useful for debugging
		})
		return
	}

	// Check if the email already exists in the database
	userExists, _ := uc.UserRepo.FindUserByEmail(req.Email)
	if userExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email already in use",
		})
		return
	}

	// Hash the password
	hashedPassword, err := utils.GenerateHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error hashing password",
			"details": err.Error(), // Include error details for hashing failure
		})
		return
	}

	// Insert the new user into the database
	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "karyawan", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.UserRepo.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Error creating user",
			"details": err.Error(), // Include the database error details
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
	})
}

func (uc *UserController) LoginUser(c *gin.Context) {
		
	//do login

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}