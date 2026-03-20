package pgsql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthRepository handles health check queries.
type HealthRepository interface {
	Check(ctx context.Context) error
}

// ─── Health ─────────────────────────────────────────────────────────────────

type healthRepository struct {
	db *pgxpool.Pool
}

func NewHealthRepository(db *pgxpool.Pool) HealthRepository {
	return &healthRepository{db: db}
}

func (r *healthRepository) Check(ctx context.Context) error {
	return r.db.Ping(ctx)
}
