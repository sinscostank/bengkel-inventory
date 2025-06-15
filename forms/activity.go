package forms

// ProductItem represents each product in the activity form
type ProductItem struct {
	ID       int `json:"id" binding:"required,gt=0"`
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

// ActivityForm ...
type ActivityForm struct {
	Products []ProductItem `json:"products" binding:"required,dive,required"`
}