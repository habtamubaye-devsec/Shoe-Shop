package database

import (
    "fmt"

    "github.com/habtamubaye-devsec/shoe-shop/internal/models"

    "gorm.io/gorm"
)

// Migrate creates and updates database tables.
func Migrate(db *gorm.DB) error {

    err := db.AutoMigrate(
        &models.User{},
        &models.Category{},
        &models.Product{},
        &models.ProductVariant{},
        &models.Cart{},
        &models.CartItem{},
        &models.Order{},
        &models.OrderItem{},
        &models.Payment{},
    )

    if err != nil {
        return fmt.Errorf(
            "database migration failed: %w",
            err,
        )
    }

    return nil
}