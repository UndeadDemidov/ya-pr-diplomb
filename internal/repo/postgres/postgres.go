package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const connTimeout = 5

// NewDB returns new Postgresql db instance.
func NewDB(c config.Postgres) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Host,
		c.Port,
		c.User,
		c.DBName,
		c.Password,
	)

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("can't connect to postgres: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("can't ping postgres: %w", err)
	}

	return db, nil
}
