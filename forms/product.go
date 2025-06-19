package forms

// RegisterForm ...
type ProductForm struct {
	Name       string  `json:"name" binding:"required"`
	Stock      int     `json:"stock" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Location   string  `json:"location" binding:"required"`
	CategoryID uint    `json:"category_id" binding:"required"`
}

type UpdateProductForm struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Location   string  `json:"location" binding:"required"`
	CategoryID uint    `json:"category_id" binding:"required"`
}