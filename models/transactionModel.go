package models

import "time"

type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	TotalAmount   uint
	Income        uint
	Expense       uint
	TransactionAt time.Time `json:"transaction_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionID"`
}

type TransactionRequest struct {
	Income        uint                         `json:"income"`
	Expense       uint                         `json:"expense"`
	TransactionAt string                       `json:"transaction_at"` // Ubah ke string
	Details       []TransactionDetailInput   `json:"details"`
}

type TransactionResponse struct {
	ID            uint                     `json:"id"`
	TotalAmount   uint                     `json:"total_amount"`
	Income        uint                     `json:"income"`
	Expense       uint                     `json:"expense"`
	TransactionAt time.Time                `json:"transaction_at"`
	Details       []TransactionDetailResponse `json:"details"`
}
