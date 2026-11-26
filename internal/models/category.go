package models

import "gorm.io/gorm"

type Category struct {
    gorm.Model

    Name        string    `gorm:"size:100;not null"`
    Slug        string    `gorm:"size:120;uniqueIndex;not null"`
    Description string    `gorm:"type:text"`

    Products []Product
}