package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	au "github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spacetab-io/pgxpoolmock"
)

var _ services.Persistent = (*User)(nil)

type User struct {
	log telemetry.AppLogger
	db  pgxpoolmock.PgxPool
}

func NewUser(database *pgxpool.Pool, logger telemetry.AppLogger) *User {
	return &User{db: database, log: logger}
}

func (r *User) Create(ctx context.Context, usr *models.User) error {
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
		insertCredStmt = `INSERT INTO gophkeeper.credentials (user_uuid, email, password) VALUES ($1, $2, $3)` //nolint:gosec
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

func (r *User) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	const selectUserQuery = `
SELECT u.uuid, crd.email, crd.password, u.created_at, u.updated_at
  FROM gophkeeper.credentials crd
  JOIN gophkeeper.users u ON u.uuid = crd.user_uuid
 WHERE crd.email=$1
`
	l := r.log.With().Str("email", email).Logger()
	usr := models.User{}
	crd := au.BasicAuth{}
	err := r.db.QueryRow(ctx, selectUserQuery, email).
		Scan(&usr.UserUUID, &crd.Email, &crd.Password, &usr.CreatedAt, &usr.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		l.Debug().Err(err).Str("query", selectUserQuery).Msg("user not found")
		return nil, fmt.Errorf("user %s not found: %w", email, err)
	}
	if err != nil {
		l.Err(err).Str("query", selectUserQuery).Msg("query error")
		return nil, fmt.Errorf("query error: %w", err)
	}
	usr.Auth = &crd
	return &usr, nil
}
