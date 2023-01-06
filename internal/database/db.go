package database

import (
	"context"
	"database/sql"
	"log"
	"time"

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

func (d *DB) CreateTable() {
	ctx, cancel := context.WithTimeout(d.ctx, 3*time.Second)
	defer cancel()

	d.DB.ExecContext(ctx, "CREATE TABLE metrics (ID SERIAL PRIMARY KEY,"+
		"NAME CHARACTER VARYING(30), "+
		"TYPE CHARACTER VARYING(10), "+
		"HASH CHARACTER VARYING(100) "+
		"VALUE DOUBLE PRECISION "+
		"DELTA INTEGER"+
		")")
}
func (d *DB) Close() {
	d.DB.Close()
}
