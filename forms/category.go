package forms

// RegisterForm ...
type CategoryForm struct {
	Name       string  `json:"name" binding:"required"`
}