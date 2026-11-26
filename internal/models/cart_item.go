package models

import "gorm.io/gorm"

type CartItem struct {
    gorm.Model

    CartID    uint `gorm:"not null;uniqueIndex:idx_cart_variant"`
    VariantID uint `gorm:"not null;uniqueIndex:idx_cart_variant"`

    Quantity int `gorm:"not null;default:1"`

    Cart    Cart
    Variant ProductVariant
}