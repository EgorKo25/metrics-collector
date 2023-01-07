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

func (d *DB) СreateTable() {
	ctx, cancel := context.WithTimeout(d.ctx, 3*time.Second)
	defer cancel()

	query := "CREATE TABLE IF NOT EXISTS metrics (id VARCHAR(30), type VARCHAR(10), hash VARCHAR(100), value DOUBLE PRECISION, delta BIGINT);"

	r, err := d.DB.ExecContext(ctx, query)
	if err != nil {
		log.Println("Field to create db table ", err, r)
	}
}
func (d *DB) Close() error {
	return d.DB.Close()
}
func (d *DB) Run(metric *storage.Metric) (err error) {

	ctx, cancel := context.WithTimeout(d.ctx, 3*time.Second)
	defer cancel()

	log.Println("! ! ! Отправка в БД ! ! ! ")
	switch metric.MType {

	case "gauge":
		r, err := d.DB.ExecContext(ctx,
			`INSERT INTO metrics 
    				(id, type, hash, value) 
					VALUES ($1, $2, $3, $4)`,
			metric.ID, metric.MType, metric.Hash, float64(*metric.Value),
		)
		if err != nil {
			log.Println("insert row into table went wrong, ", err, r)
			return err

		}

	case "counter":

		r, err := d.DB.ExecContext(ctx,
			`INSERT INTO metrics 
    				(id, type, hash, delta) 
					VALUES ($1, $2, $3, $4)`,
			metric.ID, metric.MType, metric.Hash, int(*metric.Delta),
		)
		if err != nil {
			log.Println("insert row into table went wrong, ", err, r)
			return err
		}
	}

	return nil
}
