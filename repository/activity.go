package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
	"errors"
)

// ActivityRepository defines methods to interact with the products table.
type ActivityRepository interface {
	Create(product *models.Activity) error
	FindAll() ([]models.Activity, error)
	FindByID(id uint) (*models.Activity, error)
	Update(Activity *models.Activity) error
	Delete(id uint) error
	// You can add other methods like FindByID, Update, Delete if needed
}

// ActivityRepositoryImpl is the implementation of the ActivityRepository interface.
type ActivityRepositoryImpl struct {
	DB *gorm.DB
}

// NewActivityRepository creates a new instance of ActivityRepositoryImpl
func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &ActivityRepositoryImpl{
		DB: db,
	}
}


func (r *ActivityRepositoryImpl) FindAll() ([]models.Activity, error) {
	var activities []models.Activity
	if err := r.DB.Preload("Items").Preload("User").Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *ActivityRepositoryImpl) FindByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	err := r.DB.Preload("Items").Preload("User").First(&activity, id).Error
	
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // Return nil, nil to indicate not found without error
        return nil, nil
    }

    if err != nil {
        // Real DB error
        return nil, err
    }

	return &activity, nil
}

func (r *ActivityRepositoryImpl) Create(activity *models.Activity) error {
	return r.DB.Create(activity).Error
}

func (r *ActivityRepositoryImpl) Update(activity *models.Activity) error {
	return r.DB.Save(activity).Error
}

func (r *ActivityRepositoryImpl) Delete(id uint) error {
	var activity models.Activity
	if err := r.DB.First(&activity, id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&activity).Error
}