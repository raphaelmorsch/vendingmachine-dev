package domains

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	ID int `json:"id"`

	AmountAvailable int `json:"amountAvailable"`

	Cost float32 `json:"cost"`

	ProductName string `json:"productName"`

	SellerId string `json:"sellerId"`
}
