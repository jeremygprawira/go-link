package database

import (
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"fmt"
	"time"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)
// ConnectToPostgreSQL creates a pgxpool connection pool.
func ConnectToPostgreSQL(cfg *config.Configuration) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Name,
		cfg.PostgreSQL.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	connMaxLifetime, _ := time.ParseDuration(cfg.PostgreSQL.ConnMaxLifetime)
	poolCfg.MaxConns = int32(cfg.PostgreSQL.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.PostgreSQL.MaxIdleConns)
	poolCfg.MaxConnLifetime = connMaxLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	return pool, nil
}

// DisconnectFromPostgreSQL closes the underlying connection pool.
func DisconnectFromPostgreSQL(pool *pgxpool.Pool) error {
	if pool != nil {
		pool.Close()
	}
	return nil
}
