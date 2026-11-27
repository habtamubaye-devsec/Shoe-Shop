package database

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/habtamubaye-devsec/shoe-shop/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to PostgreSQL.
func Connect(cfg *config.Config) (*gorm.DB, error) {

	// Build PostgreSQL connection URL.
	connectionURL := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			cfg.DBUser,
			cfg.DBPassword,
		),
		Host: net.JoinHostPort(
			cfg.DBHost,
			fmt.Sprintf("%d", cfg.DBPort),
		),
		Path: cfg.DBName,
	}

	// Keep the timezone slash literal for the PostgreSQL driver.
	connectionURL.RawQuery = "sslmode=disable&TimeZone=Africa/Addis_Ababa"

	// Connect using GORM.
	db, err := gorm.Open(
		postgres.Open(connectionURL.String()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to PostgreSQL: %w",
			err,
		)
	}

	// Get underlying database connection.
	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get database connection: %w",
			err,
		)
	}

	// Configure connection pool.
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Verify the database is reachable.
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()

		return nil, fmt.Errorf(
			"database ping failed: %w",
			err,
		)
	}

	return db, nil
}
