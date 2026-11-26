package routes

import (
    "net/http"

    "github.com/habtamubaye-devsec/shoe-shop/internal/handlers"
)

func RegisterRoutes(
    mux *http.ServeMux,
    productHandler *handlers.ProductHandler,
) {
    mux.HandleFunc("/api/v1/products", productHandler.Products)
    mux.HandleFunc("/api/v1/products/", productHandler.ProductByID)
}