package user

import (
	"context"
	"fmt"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	au "github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ user.Persistent = (*Repository)(nil)

type Repository struct {
	log telemetry.AppLogger
	db  *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool, logger telemetry.AppLogger) *Repository {
	return &Repository{db: database, log: logger}
}

func (r *Repository) Create(ctx context.Context, usr *models.User) error {
	auth, ok := usr.Auth.(*au.BasicAuth)
	if !ok {
		return pkg.ErrInvalidTypeCast
	}
	l := r.log.With().Object("BasicAuth", auth).Logger()
	l.Debug().Msg("insert user and credentials into db")

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction is failed: %w", err)
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				r.log.Err(err).Msg("got error on rollback transaction")

				return
			}
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				r.log.Err(err).Msg("got error on commit transaction")

				return
			}
		}
	}()

	const (
		insertUserStmt = `INSERT INTO gophkeeper.users (uuid) VALUES ($1)`
		insertCredStmt = `INSERT INTO gophkeeper.credentials (user_uuid, email, password) VALUES ($1, $2, $3)`
	)
	userUUID := uuid.New().String()
	_, err = r.db.Exec(ctx, insertUserStmt, userUUID)
	if err != nil {
		l.Err(err).Str("stmt", insertUserStmt).Msg("can't insert user into db")

		return fmt.Errorf("execution user insert statement is failed: %w", err)
	}
	_, err = r.db.Exec(ctx, insertCredStmt, userUUID, auth.Email, auth.Password)
	if err != nil {
		l.Err(err).Str("stmt", insertCredStmt).Msg("can't insert credentials into db")

		return fmt.Errorf("execution credentials insert statement is failed: %w", err)
	}
	l.Debug().Msg("user and credentials are inserted successfully")

	return nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	// TODO implement me
	panic("implement me")
}
