// activity_service.go
package service

import (
	"errors"

	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
)

type ActivityService interface {
	GetByID(id uint) (*models.Activity, error)
	Create(userID uint, userRole string, form *forms.ActivityForm) (*models.Activity, error)
	GetAll() ([]models.Activity, error)
}

type activityService struct {
	activityRepo         repository.ActivityRepository
	productRepo          repository.ProductRepository
	activityItemRepo     repository.ActivityItemRepository
	stockTransactionRepo repository.StockTransactionRepository
}

func NewActivityService(
	activityRepo repository.ActivityRepository,
	productRepo repository.ProductRepository,
	activityItemRepo repository.ActivityItemRepository,
	stockTransactionRepo repository.StockTransactionRepository,
) ActivityService {
	return &activityService{activityRepo, productRepo, activityItemRepo, stockTransactionRepo}
}

func (s *activityService) GetByID(id uint) (*models.Activity, error) {
	activity, err := s.activityRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, errors.New("activity not found")
	}
	return activity, nil
}

func (s *activityService) GetAll() ([]models.Activity, error) {
	activities, err := s.activityRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if activities == nil {
		return nil, errors.New("activity not found")
	}
	return activities, nil
}

func (s *activityService) Create(userID uint, userRole string, form *forms.ActivityForm) (*models.Activity, error) {
	// Validate activity type
	if form.Type != "inbound" && form.Type != "outbound" {
		return nil, errors.New("invalid activity type")
	}
	if form.Type == "inbound" && userRole != "admin" {
		return nil, errors.New("only admins can create inbound activities")
	}

	// Prepare maps and validate duplicates
	productIDSet := make(map[uint]struct{})
	inputMap := make(map[uint]uint)
	productIDs := make([]uint, len(form.Products))

	for i, item := range form.Products {
		if _, exists := productIDSet[item.ID]; exists {
			return nil, errors.New("duplicate product ID found")
		}
		productIDSet[item.ID] = struct{}{}
		inputMap[item.ID] = item.Quantity
		productIDs[i] = item.ID
	}

	products, err := s.productRepo.FindByIDs(productIDs)
	if err != nil || len(products) != len(form.Products) {
		return nil, errors.New("product not found or DB error")
	}

	// Check stock for outbound
	if form.Type == "outbound" {
		for _, p := range products {
			if p.Stock < int(inputMap[p.ID]) {
				return nil, errors.New("insufficient stock for product " + p.Name)
			}
		}
	}

	// Create activity
	activity := models.Activity{
		UserID: userID,
		Status: "success",
		Type:   form.Type,
	}
	if err := s.activityRepo.Create(&activity); err != nil {
		return nil, err
	}

	// Create activity items
	var activityItems []*models.ActivityItem
	for _, p := range products {
		qty := int(inputMap[p.ID])
		item := &models.ActivityItem{
			ActivityID:  activity.ID,
			ProductID:   p.ID,
			Quantity:    qty,
			PriceAtTime: p.Price,
			FinalPrice:  p.Price, // No discount
		}
		activityItems = append(activityItems, item)
	}
	if err := s.activityItemRepo.CreateMultiple(activityItems); err != nil {
		return nil, err
	}

	// Create stock transactions
	var transactions []*models.StockTransaction
	for _, item := range activityItems {
		change := item.Quantity
		if form.Type == "outbound" {
			change = -change
		}
		t := &models.StockTransaction{
			ProductID:      item.ProductID,
			ChangeQuantity: change,
			ActivityItemID: &item.ID,
			Note:           "Stock change for activity",
		}
		transactions = append(transactions, t)
	}
	if err := s.stockTransactionRepo.CreateMultiple(transactions); err != nil {
		return nil, err
	}

	// Update product stock
	for _, p := range products {
		qty := int(inputMap[p.ID])
		if form.Type == "inbound" {
			p.Stock += qty
		} else {
			p.Stock -= qty
		}
		if err := s.productRepo.Update(&p); err != nil {
			return nil, err
		}
	}

	return &activity, nil
}
