package pgsql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeremygprawira/go-link-generator/internal/models"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
)

// url_pgsql_repository.go
// UrlRepository handles url data access.
type UrlRepository interface {
	Create(ctx context.Context, input *models.Url) error
	CheckByCode(ctx context.Context, code string) (bool, error)
}

// ─── Url ────────────────────────────────────────────────────────────────────

type urlRepository struct{ db *pgxpool.Pool }

func NewUrlRepository(db *pgxpool.Pool) UrlRepository {
	return &urlRepository{db: db}
}

func (ur *urlRepository) Create(ctx context.Context, input *models.Url) error {
	logger.AddProcess(ctx, "postgre_operation", "url.create")

	_, err := ur.db.Exec(ctx,
		QueryCreateUrl,
		input.ID, input.Code, input.Name, input.Url, input.AccountNumber, input.ClickCount, input.State, input.Metadata, input.ExpiredAt,
	)

	return err
}

func (u *urlRepository) CheckByCode(ctx context.Context, code string) (bool, error) {
	logger.AddProcess(ctx, "postgre_operation", "url.check_by_code")

	var exists bool
	row := u.db.QueryRow(ctx, QueryCheckByCode, code)
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
