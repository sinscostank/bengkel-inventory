package service

import (
	"errors"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
)

type CategoryService interface {
	GetAll(page, limit int) ([]models.Category, int64, error)
	GetByID(id uint) (*models.Category, error)
	Create(form *forms.CategoryForm) (*models.Category, error)
	Update(id uint, form *forms.CategoryForm) (*models.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo}
}

func (s *categoryService) GetAll(page, limit int) ([]models.Category, int64, error) {
	return s.categoryRepo.FindAll(page, limit)
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) Create(form *forms.CategoryForm) (*models.Category, error) {
	category := &models.Category{
		Name: form.Name,
	}
	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Update(id uint, form *forms.CategoryForm) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	category.Name = form.Name
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Delete(id uint) error {
	return s.categoryRepo.Delete(id)
}
