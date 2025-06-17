// repository/user_repository.go
package repository

import (
	"errors"

	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
)


type UserRepository interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

// ProductRepositoryImpl is the implementation of the ProductRepository interface.
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepository returns a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

// FindByEmail retrieves a user by their email
func (r *UserRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if the user is not found
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user in the database
func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}
