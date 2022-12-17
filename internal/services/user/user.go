package user

import (
	"context"
	"fmt"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	au "github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
)

type Persistent interface {
	Create(context.Context, *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type Service struct {
	// ToDo add authenticator base on context.
	// https://github.com/johanbrandhorst/grpc-auth-example
	persist Persistent
}

func NewService(repository Persistent) *Service {
	return &Service{persist: repository}
}

func (s *Service) SignUp(ctx context.Context, usr *models.User) error {
	auth, ok := usr.Auth.(*au.BasicAuth)
	if !ok {
		return pkg.ErrInvalidTypeCast
	}

	existsUser, err := s.findByEmail(ctx, auth)
	if existsUser != nil || err == nil {
		return pkg.ErrEmailExists
	}
	err = auth.HashPassword()
	if err != nil {
		return fmt.Errorf("failed with hash password: %w", err)
	}

	return s.persist.Create(ctx, usr) //nolint:wrapcheck
}

func (s *Service) SignIn(ctx context.Context, auth *au.BasicAuth) (*models.User, error) {
	foundUser, err := s.findByEmail(ctx, auth)
	if err != nil {
		return nil, fmt.Errorf("given email %s not found: %w", auth.Email, err)
	}

	err = foundUser.Auth.(*au.BasicAuth).ValidatePassword(auth.Password) //nolint:forcetypeassert
	if err != nil {
		return nil, fmt.Errorf("given password is not valid: %w", err)
	}

	return foundUser, nil
}

func (s *Service) findByEmail(ctx context.Context, auth *au.BasicAuth) (*models.User, error) {
	auth.CleanCredentials() // не явная чистка значений
	foundUser, err := s.persist.FindByEmail(ctx, auth.Email)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return foundUser, nil
}
