package go_postgres

import (
	"context"
)

func NewTestClient() (*PostgresClient, error) {
	cfg := &IPGConfig{
		DSN:    "postgres://admin:8246@localhost:5433/aggregation_service_db?sslmode=disable",
		Driver: "postgres",
	}
	return NewPGClient(context.Background(), *cfg)
}
