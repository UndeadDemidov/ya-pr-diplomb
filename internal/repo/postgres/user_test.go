package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/telemetry"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spacetab-io/pgxpoolmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	type fields struct {
		pool *pgxpoolmock.MockPgxIface
		tx   pgx.Tx
	}
	type args struct {
		usr *models.User
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "inserted successfully",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil),
					f.pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil),
					f.pool.EXPECT().Commit(gomock.Any()).Return(nil),
				)
			},
			args:    args{models.NewUser(auth.NewBasicAuth("test", "test"))},
			wantErr: false,
		},
		{
			name: "user inserted unsuccessfully",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, pkg.ErrDumb),
					f.pool.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
			args:    args{models.NewUser(auth.NewBasicAuth("test", "test"))},
			wantErr: true,
		},
		{
			name: "credentials inserted unsuccessfully",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil),
					f.pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, pkg.ErrDumb),
					f.pool.EXPECT().Rollback(gomock.Any()).Return(nil),
				)
			},
			args:    args{models.NewUser(auth.NewBasicAuth("test", "test"))},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

			// begin tx - given
			mockPool.EXPECT().Begin(gomock.Any()).Return(mockPool, nil).AnyTimes()
			tx, err := mockPool.Begin(context.Background())
			assert.NoError(t, err)

			if tt.prepare != nil {
				tt.prepare(&fields{pool: mockPool, tx: tx})
			}

			r := &User{
				log: telemetry.NewTestAppLogger(),
				db:  mockPool,
			}
			// when - then
			if err = r.Create(context.Background(), tt.args.usr); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_FindByEmail(t *testing.T) {
	type (
		fields struct {
			pool *pgxpoolmock.MockPgxIface
		}
		args struct {
			email string
		}
	)
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "user returned successfully",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(
						pgxpoolmock.NewRow(
							uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
							"test", "test",
							time.Unix(1671301348, 0),
							time.Unix(1671301348, 0),
						)),
				)
			},
			args: args{email: "test"},
			want: &models.User{
				UserUUID:  uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
				Auth:      auth.NewBasicAuth("test", "test"),
				CreatedAt: time.Unix(1671301348, 0),
				UpdatedAt: time.Unix(1671301348, 0),
			},
			wantErr: assert.NoError,
		},
		{
			name: "query error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(
						pgxpoolmock.NewRow(nil).WithError(pkg.ErrDumb)))
			},
			args:    args{email: "test"},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "user not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(
						pgxpoolmock.NewRow(nil).WithError(sql.ErrNoRows)))
			},
			args:    args{email: "test"},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			// given
			if tt.prepare != nil {
				tt.prepare(&fields{pool: mockPool})
			}

			r := &User{
				log: telemetry.NewTestAppLogger(),
				db:  mockPool,
			}

			// when - then
			got, err := r.FindByEmail(context.Background(), tt.args.email)
			if !tt.wantErr(t, err, fmt.Sprintf("FindByEmail(ctx, %v)", tt.args.email)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindByEmail(ctx, %v)", tt.args.email)
		})
	}
}
