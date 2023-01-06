package database

import (
	"database/sql"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
)

type DB struct {
	DB *sql.DB
}

func NewDB(cfg *config.ConfigurationServer) *DB {
	if cfg.DB == "" {
		return nil
	}
	db, err := sql.Open("pgx",
		cfg.DB)
	if err != nil {
		panic(err)
	}
	return &DB{
		DB: db,
	}
}

func (d *DB) Close() {
	d.DB.Close()
}
