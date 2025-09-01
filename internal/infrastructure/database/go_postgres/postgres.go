package go_postgres

import (
	"AggregationService/internal/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type PostgresClient struct {
	DB      *sqlx.DB
	Builder squirrel.StatementBuilderType
}

func NewPGClient(ctx context.Context, cfg IPGConfig) (*PostgresClient, error) {
	db, err := sqlx.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with db: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	logger.FromContext(ctx).
		Debug(fmt.Sprintf("pinging database with dsn: %w", cfg.DSN))
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	logger.FromContext(ctx).
		Debug(fmt.Sprintf("successfully connected to db: %w", cfg.DSN))
	return &PostgresClient{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (p *PostgresClient) SQLDB() *sql.DB {
	return p.DB.DB
}
