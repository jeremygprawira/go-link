package pgsql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeremygprawira/go-link-generator/internal/models"
)

// user_pgsql_repository.go
// User-specific repository methods are implemented in health_pgsql_repository.go.
// Add additional query-specific methods here as your project grows.
// UserRepository handles user data access.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	CheckByEmailOrPhoneNumber(ctx context.Context, email string, phoneNumber string) (bool, error)
	GetOneByAccountNumber(ctx context.Context, accountNumber string) (*models.User, error)
	GetCredentialsByEmailOrPhoneNumber(ctx context.Context, email string, phoneNumber string) (*models.User, error)
}

// ─── User ────────────────────────────────────────────────────────────────────

type userRepository struct{ db *pgxpool.Pool }

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := ur.db.Exec(ctx,
		`INSERT INTO users (id, account_number, name, email, phone_number, phone_country_code, password)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.AccountNumber, user.Name, user.Email, user.PhoneNumber, user.PhoneCountryCode, user.Password,
	)
	return err
}

func (ur *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	row := ur.db.QueryRow(ctx,
		`SELECT id, account_number, name, email, phone_number, phone_country_code, created_at, updated_at FROM users WHERE id = $1`,
		id,
	)
	if err := row.Scan(&user.ID, &user.AccountNumber, &user.Name, &user.Email, &user.PhoneNumber, &user.PhoneCountryCode, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) CheckByEmailOrPhoneNumber(ctx context.Context, email string, phoneNumber string) (bool, error) {
	var exists bool
	row := ur.db.QueryRow(ctx, QueryCheckByEmailOrPhoneNumber, email, phoneNumber)
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (ur *userRepository) GetOneByAccountNumber(ctx context.Context, accountNumber string) (*models.User, error) {
	var user models.User
	row := ur.db.QueryRow(ctx, QueryGetByAccountNumber, accountNumber)
	if err := row.Scan(&user.ID, &user.AccountNumber, &user.Name, &user.Email, &user.PhoneNumber, &user.PhoneCountryCode, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetCredentialsByEmailOrPhoneNumber(ctx context.Context, email string, phoneNumber string) (*models.User, error) {
	var user models.User
	row := ur.db.QueryRow(ctx, QueryGetCredentialsByEmailOrPhoneNumber, email, phoneNumber)
	if err := row.Scan(&user.ID, &user.AccountNumber, &user.Name, &user.Email, &user.PhoneNumber, &user.PhoneCountryCode, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}
