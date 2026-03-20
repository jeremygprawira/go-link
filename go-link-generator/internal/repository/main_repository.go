package repository

import (
	"github.com/jeremygprawira/go-link-generator/internal/pkg/database"
	"github.com/jeremygprawira/go-link-generator/internal/repository/pgsql"
)

// Repository aggregates all data-layer repositories.
type Repository struct {
	Postgre *pgsql.PostgreRepository
}

func New(db *database.Database) *Repository {
	return &Repository{
		Postgre: pgsql.New(db.PostgreDatabase),
	}
}
