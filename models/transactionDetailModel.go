package models

import "time"

type TransactionDetail struct {
	ID            uint    `gorm:"primaryKey"`
	TransactionID uint    `json:"transaction_id"`
	ProductID     uint    `json:"product_id"`
	Product       Product `gorm:"foreignKey:ProductID"`
	Quantity      uint    `json:"quantity"`
	SubTotal      uint    `json:"sub_total"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionDetailInput struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type TransactionDetailResponse struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
	SubTotal  uint `json:"sub_total"`
}
