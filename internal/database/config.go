package database

import (
	"database/sql"
	"fmt"
	"garage-api/internal/config"
	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
} 