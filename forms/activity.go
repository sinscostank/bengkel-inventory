package forms

// ProductItem represents each product in the activity form
type ProductItem struct {
	ID       uint `json:"id" binding:"required,gt=0"`
	Quantity uint `json:"quantity" binding:"required,gt=0"`
}

// ActivityForm ...
type ActivityForm struct {
	Products []ProductItem `json:"products" binding:"required,dive,required"`
	Type	 string        `json:"type" binding:"required,oneof=outbound inbound"`
}