package user

import (
	"context"
	"fmt"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	"github.com/opentracing/opentracing-go"
)

type Persistent interface {
	Create(context.Context, *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type Service struct {
	persist Persistent
}

func NewService(repository Persistent) *Service {
	return &Service{persist: repository}
}

func (s *Service) SignIn(ctx context.Context, usr *models.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.SignIn")
	defer span.Finish()

	auth, ok := usr.Auth.(*pkg.BasicAuth)
	if !ok {
		return pkg.ErrInvalidTypeCast
	}

	existsUser, err := s.findByEmail(ctx, auth)
	if existsUser != nil || err == nil {
		return pkg.ErrEmailExists
	}

	return s.persist.Create(ctx, usr) //nolint:wrapcheck
}

func (s *Service) SignOn(ctx context.Context, auth *pkg.BasicAuth) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Service.SignOn")
	defer span.Finish()

	foundUser, err := s.findByEmail(ctx, auth)
	if err != nil {
		return nil, fmt.Errorf("given email %s not found: %w", auth.Email, err)
	}

	err = foundUser.Auth.(*pkg.BasicAuth).ValidatePassword(auth.Password) //nolint:forcetypeassert
	if err != nil {
		return nil, fmt.Errorf("given password is not valid: %w", err)
	}

	foundUser.Sanitize()
	return foundUser, nil
}

func (s *Service) findByEmail(ctx context.Context, auth *pkg.BasicAuth) (*models.User, error) {
	auth.CleanCredentials()
	foundUser, err := s.persist.FindByEmail(ctx, auth.Email)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return foundUser, nil
}
