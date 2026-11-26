package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	UserID uint `gorm:"uniqueIndex;not null"`

	Items []CartItem

	User *User
}
