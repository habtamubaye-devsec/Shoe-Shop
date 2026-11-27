package models

import "gorm.io/gorm"

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model

	OrderNumber string `gorm:"size:50;uniqueIndex;not null"`

	UserID uint `gorm:"not null;index"`

	Status OrderStatus `gorm:"size:30;not null;default:pending"`

	Currency string `gorm:"size:3;not null;default:ETB"`

	SubtotalMinor int64 `gorm:"not null"`
	TotalMinor    int64 `gorm:"not null"`

	Items    []OrderItem
	Payments []Payment

	User User
}
