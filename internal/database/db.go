package database

import (
	"context"
	"log"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	DB  *pgx.Conn
	ctx context.Context
}

func NewDB(cfg *config.ConfigurationServer, ctx context.Context) *DB {
	if cfg.DB == "" {
		return nil
	}

	conn, err := pgx.Connect(ctx, cfg.DB)
	if err != nil {
		log.Println(err)
	}
	return &DB{
		DB:  conn,
		ctx: ctx,
	}
}

func (d *DB) Close() {
	d.DB.Close(d.ctx)
}
