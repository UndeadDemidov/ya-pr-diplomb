package services

import (
	"context"
	"testing"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services/mocks"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestService_SignUp(t *testing.T) {
	type fields struct {
		persist *mock_services.MockPersistent
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
			name: "signup successfully, user not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.persist.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, nil),
					f.persist.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			args: args{
				&models.User{
					UserUUID:  uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
					Auth:      auth.NewBasicAuth("test", "test"),
					CreatedAt: time.Unix(1671301348, 0),
					UpdatedAt: time.Unix(1671301348, 0),
				},
			},
			wantErr: false,
		},
		{
			name: "signup successfully, got persist error",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.persist.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(nil, pkg.ErrDumb),
					f.persist.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			args: args{
				&models.User{
					UserUUID:  uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
					Auth:      auth.NewBasicAuth("test", "test"),
					CreatedAt: time.Unix(1671301348, 0),
					UpdatedAt: time.Unix(1671301348, 0),
				},
			},
			wantErr: false,
		},
		{
			name: "signup failed, user found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.persist.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(&models.User{}, nil),
				)
			},
			args: args{
				&models.User{
					UserUUID:  uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
					Auth:      auth.NewBasicAuth("test", "test"),
					CreatedAt: time.Unix(1671301348, 0),
					UpdatedAt: time.Unix(1671301348, 0),
				},
			},
			wantErr: true,
		},
		{
			name:    "failed au cast",
			prepare: nil,
			args: args{
				&models.User{
					UserUUID:  uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
					Auth:      nil,
					CreatedAt: time.Unix(1671301348, 0),
					UpdatedAt: time.Unix(1671301348, 0),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			p := mock_services.NewMockPersistent(ctrl)
			f := fields{
				persist: p,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			s := NewUser(p)

			if err := s.SignUp(context.Background(), tt.args.usr); (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_SignIn(t *testing.T) {
	type fields struct {
		persist *mock_services.MockPersistent
		usr     *models.User
	}
	type args struct {
		au *auth.BasicAuth
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		found   *models.User
		wantErr bool
	}{
		{
			name: "signin successfully",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.persist.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(f.usr, nil),
				)
			},
			args: args{auth.NewBasicAuth("test", "test")},
			found: &models.User{
				UserUUID:  uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
				Auth:      auth.NewBasicAuth("test", "test"),
				CreatedAt: time.Unix(1671301348, 0),
				UpdatedAt: time.Unix(1671301348, 0),
			},
			wantErr: false,
		},
		{
			name: "signin failed, user not found",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.persist.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(f.usr, pkg.ErrDumb),
				)
			},
			args:    args{auth.NewBasicAuth("test", "test")},
			found:   nil,
			wantErr: true,
		},
		{
			name: "signin failed, password not valid",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.persist.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(f.usr, pkg.ErrDumb),
				)
			},
			args: args{auth.NewBasicAuth("test", "invalid")},
			found: &models.User{
				UserUUID:  uuid.MustParse("0ad66d2e-fc9e-4c16-8355-a6b9f89866d7"),
				Auth:      auth.NewBasicAuth("test", "test"),
				CreatedAt: time.Unix(1671301348, 0),
				UpdatedAt: time.Unix(1671301348, 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			p := mock_services.NewMockPersistent(ctrl)
			f := fields{
				persist: p,
				usr:     tt.found,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			s := NewUser(p)

			if tt.found != nil {
				err := tt.found.Auth.(*auth.BasicAuth).HashPassword()
				require.NoError(t, err)
			}

			_, err := s.SignIn(context.Background(), tt.args.au)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
