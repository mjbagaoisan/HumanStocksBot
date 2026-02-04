package repos

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		t.Skip("Skipping integration test: no database URL configured")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("Skipping integration test: %v", err)
	}
	
	return pool
}
