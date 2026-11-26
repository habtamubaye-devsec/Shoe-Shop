package models

import "gorm.io/gorm"

type OrderItem struct {
    gorm.Model

    OrderID uint `gorm:"not null;index"`

    VariantID *uint `gorm:"index"`

    ProductName string `gorm:"size:150;not null"`
    SKU         string `gorm:"size:100;not null"`

    Size  uint   `gorm:"not null"`
    Color string `gorm:"size:50;not null"`

    Quantity int `gorm:"not null"`

    UnitPriceMinor int64 `gorm:"not null"`
    LineTotalMinor int64 `gorm:"not null"`

    Order   Order
    Variant *ProductVariant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}