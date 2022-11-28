package user

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ user.Persistent = (*Repository)(nil)

type Repository struct {
	db *pgxpool.Pool
}

func (r *Repository) Create(ctx context.Context, m *models.User) error {
	// TODO implement me
	panic("implement me")
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	// TODO implement me
	panic("implement me")
}
