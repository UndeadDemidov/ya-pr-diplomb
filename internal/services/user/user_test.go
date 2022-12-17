package user

import (
	"context"
	"testing"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	mock_user "github.com/UndeadDemidov/ya-pr-diplomb/internal/services/user/mocks"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestService_SignUp(t *testing.T) {
	type fields struct {
		persist *mock_user.MockPersistent
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
			name:    "failed auth cast",
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

			p := mock_user.NewMockPersistent(ctrl)
			f := fields{
				persist: p,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			s := NewService(p)

			if err := s.SignUp(context.Background(), tt.args.usr); (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
