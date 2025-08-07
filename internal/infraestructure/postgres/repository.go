package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	// Import the pgx driver
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/renzonaitor/tweet-api/cmd/http/config"
)

// Repository holds the database connection pool.
type Repository struct {
	db *sql.DB
}

func NewRepository(cfg config.Config) *Repository {
	// 1. Construct the Data Source Name (DSN) string from your config.
	// This string contains all the necessary info to connect to the database.
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)

	// 2. Open a connection pool.
	// `sql.Open` doesn't actually create any connections yet, it just prepares the pool.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		// If this fails, the application can't start, so we panic.
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// 3. Configure the connection pool.
	// These are good default settings for a web service.
	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConnection) // Max number of open connections
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConnection) // Max number of connections in the idle pool
	db.SetConnMaxLifetime(5 * time.Minute)             // Max time a connection can be reused

	// 4. Verify the connection is alive.
	// `db.Ping` is crucial to ensure the database is reachable and credentials are correct.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the PostgresSQL database.")

	// 5. Return the repository with the active connection pool.
	return &Repository{
		db: db,
	}
}

// Close gracefully closes the database connection pool.
func (r *Repository) Close() {
	if err := r.db.Close(); err != nil {
		log.Printf("Error closing the database: %v", err)
	}
}
