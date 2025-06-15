package repository

import (
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/gorm"
)

// StockTransactionRepository defines methods to interact with the products table.
type StockTransactionRepository interface {
	Create(product *models.StockTransaction) error
	CreateStockTransactions(StockTransactions []*models.StockTransaction) error
}

// StockTransactionRepositoryImpl is the implementation of the StockTransactionRepository interface.
type StockTransactionRepositoryImpl struct {
	DB *gorm.DB
}

// NewStockTransactionRepository creates a new instance of StockTransactionRepositoryImpl
func NewStockTransactionRepository(db *gorm.DB) StockTransactionRepository {
	return &StockTransactionRepositoryImpl{
		DB: db,
	}
}

func (r *StockTransactionRepositoryImpl) Create(StockTransaction *models.StockTransaction) error {
	return r.DB.Create(StockTransaction).Error
}

func (r *StockTransactionRepositoryImpl) CreateStockTransactions(StockTransactions []*models.StockTransaction) error {
	if len(StockTransactions) == 0 {
		return nil // No items to create
	}

	return r.DB.Create(StockTransactions).Error
}