package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Category string
	Price    float64
	SKU      string
	Stock    int
}
