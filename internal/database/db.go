// Package database пакет организует работу с базой данных
// Работа организована через объект структуры DB
//
// Объект создаеться при помощи команды:
//
//	newDB := NewDB(config, storage)
package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/EgorKo25/DevOps-Track-Yandex/internal/configuration"
	"github.com/EgorKo25/DevOps-Track-Yandex/internal/storage"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB структура для работы с базой данных
type DB struct {
	DB  *sql.DB
	cfg *config.ConfigurationServer
	str *storage.MetricStorage

	Buffer []storage.Metric
}

// NewDB конструктор для создания экземпяра структуры DB
func NewDB(cfg *config.ConfigurationServer, str *storage.MetricStorage) *DB {

	ctx := context.Background()

	if cfg.DB == "" {
		return nil
	}

	db, err := sql.Open("pgx", cfg.DB)
	if err != nil {
		log.Println(err)
	}

	createTableWithContext(ctx, db)

	return &DB{
		DB:     db,
		cfg:    cfg,
		str:    str,
		Buffer: make([]storage.Metric, 0, 29),
	}
}

// createTableWithContext создаёт все необходимые таблицы в базе данных
func createTableWithContext(ctx context.Context, db *sql.DB) {

	childCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := "CREATE TABLE IF NOT EXISTS metrics (id VARCHAR(30), type VARCHAR(10), hash VARCHAR(100), value DOUBLE PRECISION, delta BIGINT);"

	r, err := db.ExecContext(childCtx, query)
	if err != nil {
		log.Println("Field to create db table ", err, r)
	}
}

// Close закрывает соединение с базой данных
func (d *DB) Close() error {
	return d.DB.Close()
}

// FlushWithContext отправляет транзакицю в базу данных
func (d *DB) FlushWithContext(ctx context.Context) (err error) {

	childCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if d.DB == nil {
		return errors.New("you haven`t opened the database connection")
	}

	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}

	var stmt *sql.Stmt
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println("Statement error: ", err)
		}
	}()

	for _, v := range d.Buffer {

		switch v.MType {
		case "gauge":
			query := `INSERT INTO metrics 
    				(id, type, hash, value) 
					VALUES ($1, $2, $3, $4)`

			stmt, err = tx.PrepareContext(childCtx, query)
			if err != nil {
				return err
			}

			if _, err = stmt.ExecContext(childCtx, v.ID, v.MType, v.Hash, float64(*v.Value)); err != nil {
				if err = tx.Rollback(); err != nil {
					log.Println("update drivers: unable to rollback: ", err)
				}
				return err
			}

		case "counter":
			query := `INSERT INTO metrics 
    				(id, type, hash, delta) 
					VALUES ($1, $2, $3, $4)`

			stmt, err = tx.PrepareContext(childCtx, query)
			if err != nil {
				return err
			}

			if _, err = stmt.ExecContext(childCtx, v.ID, v.MType, v.Hash, int(*v.Delta)); err != nil {
				if err = tx.Rollback(); err != nil {
					log.Println("update drivers: unable to rollback: ", err)
				}
				return err
			}

		}

	}

	if err = tx.Commit(); err != nil {
		log.Println("update drivers: unable to commit: ", err)
		return err
	}

	d.Buffer = d.Buffer[:0]
	return nil

}

// Run добавляет метрику в базу данных
func (d *DB) Run(ctx context.Context, metric *storage.Metric) (err error) {

	childCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	switch metric.MType {

	case "gauge":
		_, err = d.DB.ExecContext(childCtx,
			`INSERT INTO metrics 
    				(id, type, hash, value) 
					VALUES ($1, $2, $3, $4)`,
			metric.ID, metric.MType, metric.Hash, float64(*metric.Value),
		)
		if err != nil {
			log.Println("insert row into table went wrong, ", err)
			return err

		}

	case "counter":

		_, err = d.DB.ExecContext(childCtx,
			`INSERT INTO metrics 
    				(id, type, hash, delta) 
					VALUES ($1, $2, $3, $4)`,
			metric.ID, metric.MType, metric.Hash, int(*metric.Delta),
		)
		if err != nil {
			log.Println("insert row into table went wrong, ", err)
			return err
		}
	}

	return nil
}
