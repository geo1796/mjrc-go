package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"mjrc/core/postgres/dao"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB interface {
	PingWithRetry(ctx context.Context, retries int, delay time.Duration) error
	Pool() *pgxpool.Pool
	Queries() *dao.Queries
	Begin(ctx context.Context) (pgx.Tx, error)
	Migrate() error
	Close()
}

type db struct {
	pool    *pgxpool.Pool
	queries *dao.Queries
}

func (db *db) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *db) Queries() *dao.Queries {
	return db.queries
}

func (db *db) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

func (db *db) Migrate() error {
	sqlDB := stdlib.OpenDBFromPool(db.Pool())
	goose.SetBaseFS(migrationsFS)

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("failed to migrate db: %v", err)
	}
	return nil
}

func (db *db) Close() {
	db.pool.Close()
}

func New(ctx context.Context, dsn string,
	connMaxLifetime, connMaxIdleTime time.Duration,
	maxOpenConns, maxIdleConns int) (DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %v", err)
	}

	cfg.MaxConnLifetime = connMaxLifetime
	cfg.MaxConnIdleTime = connMaxIdleTime
	cfg.MaxConns = int32(maxOpenConns)
	cfg.MinConns = int32(maxIdleConns)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %v", err)
	}

	return &db{
		pool:    pool,
		queries: dao.New(pool),
	}, nil
}

func (db *db) PingWithRetry(ctx context.Context, retries int, delay time.Duration) error {
	if retries < 1 {
		return errors.New("pingWithRetry: retries must be >= 1")
	}
	pool := db.Pool()
	var err error
	for i := range retries {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err = pool.Ping(pingCtx)
		cancel()
		if err == nil {
			return nil
		}
		if i == retries-1 {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		if delay < time.Second {
			delay *= 2 // backoff plafonné implicitement
		}
	}
	return err
}
