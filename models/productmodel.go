package models

import "time"

type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Stock       uint
	Price       uint
	UnitsSold   uint
	ImageUrl    string
	CategoryID  uint
	Category    Category `gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Stock       *uint   `json:"stock"`
	Price       *uint   `json:"price"`
	UnitsSold	*uint   `json:"units_sold"`
	ImageUrl    *string `json:"image_url"`
	CategoryID  *uint   `json:"category_id"`
}

type ProductResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Stock        uint      `json:"stock"`
	Price        uint      `json:"price"`
	UnitsSold    uint      `json:"units_sold"`
	ImageUrl     string    `json:"image_url"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
