package models

import "gorm.io/gorm"

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	gorm.Model

	FullName     string   `gorm:"size:150;not null"`
	Email        string   `gorm:"size:150;uniqueIndex;not null"`
	PasswordHash string   `gorm:"not null"`
	Phone        string   `gorm:"size:20"`
	Role         UserRole `gorm:"size:20;not null;default:customer"`
	IsActive     bool     `gorm:"not null;default:true"`

	Cart   *Cart
	Orders []Order
}
