package models

import "gorm.io/gorm"

type Product struct {
    gorm.Model

    Name        string `gorm:"size:150;not null"`
    Slug        string `gorm:"size:180;uniqueIndex;not null"`
    Brand       string `gorm:"size:100"`
    Description string `gorm:"type:text"`

    // Price is stored in Ethiopian santim.
    PriceMinor int64 `gorm:"not null;default:0"`

    CategoryID uint     `gorm:"not null;index"`
    Category   Category

    Variants []ProductVariant

    IsActive bool `gorm:"not null;default:true"`
}