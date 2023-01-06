package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	DB  *sql.DB
	ctx context.Context
}

func NewDB(cfg *config.ConfigurationServer, ctx context.Context) *DB {
	if cfg.DB == "" {
		return nil
	}

	conn, err := sql.Open("pgx", cfg.DB)
	if err != nil {
		log.Println(err)
	}
	return &DB{
		DB:  conn,
		ctx: ctx,
	}
}

func (d *DB) Close() {
	d.DB.Close()
}
