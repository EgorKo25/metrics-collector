package database

import (
	"context"
	"database/sql"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"
	"log"
	"time"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	DB  *sql.DB
	ctx context.Context
	cfg *config.ConfigurationServer
	str *storage.MetricStorage
}

func NewDB(cfg *config.ConfigurationServer, ctx context.Context, str *storage.MetricStorage) *DB {
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
		cfg: cfg,
		str: str,
	}
}

func (d *DB) CreateTable() {
	ctx, cancel := context.WithTimeout(d.ctx, 3*time.Second)
	defer cancel()

	d.DB.ExecContext(ctx, `CREATE TABLE metrics (ID PRIMARY KEY, 
		name varchar(30), 
		type varchar(10), 
		hash varchar(100), 
		value double precision, 
		delta integer
		);`)
}
func (d *DB) Close() error {
	return d.DB.Close()
}
func (d *DB) WriteAll() (err error) {
	var metric storage.Metric

	ctx, cancel := context.WithTimeout(d.ctx, 3*time.Second)
	defer cancel()

	for k, v := range d.str.Metrics {

		metric.ID = k
		metric.MType = v.MType

		if v.MType == "gauge" {
			metric.Value = v.Value
			metric.Delta = nil
		}

		if v.MType == "counter" {
			metric.Value = nil
			metric.Delta = v.Delta
		}

		d.DB.ExecContext(ctx,
			`INSERT INTO metrics (NAME, TYPE, HASH, VALUE, DELTA) VALUES (@name, @type, @hash, @value, @delta)`,
			sql.Named("name", metric.ID),
			sql.Named("hash", metric.Hash),
			sql.Named("type", metric.MType),
			sql.Named("value", metric.Value),
			sql.Named("delta", metric.Delta),
		)

	}

	return d.Close()
}

func (d *DB) Run() error {
	tickerSave := time.NewTicker(d.cfg.StoreInterval)

	for {
		select {
		case <-tickerSave.C:
			if err := d.WriteAll(); err != nil {
				return err
			}
		}
	}
}
