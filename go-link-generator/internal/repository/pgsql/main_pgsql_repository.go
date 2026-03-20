package pgsql

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreRepository aggregates all PostgreSQL repositories.
type PostgreRepository struct {
	Health HealthRepository
	User   UserRepository
}

func New(db *pgxpool.Pool) *PostgreRepository {
	return &PostgreRepository{
		Health: NewHealthRepository(db),
		User:   NewUserRepository(db),
	}
}
