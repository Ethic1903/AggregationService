package migrations

import (
	"embed"
	"fmt"

	"AggregationService/internal/infrastructure/database/go_postgres"
	"github.com/pressly/goose/v3"
)

var (
	//go:embed subscription/*.sql
	migrations embed.FS
)

func MigrateDB(db *go_postgres.PostgresClient) error {
	if err := migrate(db, "subscription"); err != nil {
		return fmt.Errorf("migrate: %v", err)
	}
	return nil
}

func migrate(db *go_postgres.PostgresClient, dir string) error {
	goose.SetBaseFS(migrations)
	if err := goose.Up(db.SQLDB(), dir); err != nil {
		return fmt.Errorf("goose up: %v", err)
	}
	return nil
}
