package models

import "gorm.io/gorm"

type PaymentStatus string

const (
    PaymentPending PaymentStatus = "pending"
    PaymentSuccess PaymentStatus = "success"
    PaymentFailed  PaymentStatus = "failed"
    PaymentRefunded PaymentStatus = "refunded"
)

type Payment struct {
    gorm.Model

    OrderID uint `gorm:"not null;index"`

    AmountMinor int64 `gorm:"not null"`

    Currency string `gorm:"size:3;not null;default:ETB"`

    Provider          string `gorm:"size:50"`
    ProviderReference string `gorm:"size:150;index"`

    Status PaymentStatus `gorm:"size:30;not null;default:pending"`

    Order Order
}