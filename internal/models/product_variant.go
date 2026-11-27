package models

import "gorm.io/gorm"

type ProductVariant struct {
	gorm.Model

	ProductID uint `gorm:"not null;uniqueIndex:idx_product_size_color,priority:1"`

	SKU   string `gorm:"size:100;uniqueIndex;not null"`
	Size  uint   `gorm:"not null;uniqueIndex:idx_product_size_color,priority:2"`
	Color string `gorm:"size:50;not null;uniqueIndex:idx_product_size_color,priority:3"`

	Stock int `gorm:"not null;default:0;check:stock >= 0"`

	PriceMinor *int64

	Product Product
}
