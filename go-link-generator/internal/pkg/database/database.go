package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeremygprawira/go-link-generator/internal/config"
)

// Database holds all connected database clients.
type Database struct {
	PostgreDatabase *pgxpool.Pool
}

// Connect initialises all selected databases and returns the aggregate struct.
func Connect(config *config.Configuration) (*Database, error) {
	postgreSQL, err := ConnectToPostgreSQL(config)
	if err != nil {
		return nil, err
	}

	return &Database{
		PostgreDatabase: postgreSQL,
	}, nil
}

func Disconnect(db *Database) error {
	return DisconnectFromPostgreSQL(db.PostgreDatabase)
}
