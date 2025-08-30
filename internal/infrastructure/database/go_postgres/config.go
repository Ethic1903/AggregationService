package go_postgres

import (
	"fmt"
	"os"
)

type IPGConfig struct {
	DSN    string
	Driver string
}

func NewConfig() (*Config, error) {
	dsn := os.Getenv("PG_DSN")
	driver := os.Getenv("PG_DRIVER")
	if dsn == "" || driver == "" {
		return nil, fmt.Errorf("could not get dsn/driver for db")
	}
	return &Config{
		DSN:    dsn,
		Driver: driver,
	}, nil
}
