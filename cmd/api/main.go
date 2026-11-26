package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/habtamubaye-devsec/shoe-shop/internal/config"
	"github.com/habtamubaye-devsec/shoe-shop/internal/database"
	"github.com/habtamubaye-devsec/shoe-shop/internal/handlers"
	"github.com/habtamubaye-devsec/shoe-shop/internal/repositories"
	"github.com/habtamubaye-devsec/shoe-shop/internal/routes"
	"github.com/habtamubaye-devsec/shoe-shop/internal/services"
)

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Load application configuration.
	cfg, err := config.LoadConfig()

	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Connect to PostgreSQL.
	db, err := database.Connect(cfg)

	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	// Run database migrations.
	if err := database.Migrate(db); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	log.Println("Database migrations completed successfully")

	// Get underlying SQL connection.
	sqlDB, err := db.DB()

	if err != nil {
		return fmt.Errorf("failed to get SQL connection: %w", err)
	}

	// Close database connection when application exits.
	defer sqlDB.Close()

	// Dependency injection: build each layer and pass it to the next.
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize HTTP router.
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, productHandler)

	// Register health endpoint.
	mux.HandleFunc("/api/v1/health", healthCheck)

	// Create HTTP server.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: mux,
	}

	// Application information.
	log.Println("--------------------------------")
	log.Println("Application:", cfg.AppName)
	log.Println("Environment:", cfg.AppEnv)
	log.Println("Port:", cfg.AppPort)
	log.Println("Database: Connected successfully")
	log.Println("--------------------------------")

	// Start HTTP server.
	err = server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return fmt.Errorf("HTTP server error: %w", err)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, `{
        "success": true,
        "message": "Shoe Shop API is running"
    }`)
}
