package models

type ProductSales struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Stock      int    `json:"stock"`
	TotalSales int    `json:"total_sales"`
}