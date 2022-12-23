package services

import (
	"context"
	"fmt"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/delivery"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	au "github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	_ "github.com/golang/mock/mockgen/model" //
)

//go:generate mockgen -destination=./mocks/mock_user.go . Persistent

type Persistent interface {
	Create(context.Context, *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

var _ delivery.User = (*User)(nil)

type User struct {
	// ToDo add authenticator base on context.
	// https://github.com/johanbrandhorst/grpc-auth-example
	persist Persistent
}

func NewUser(repository Persistent) *User {
	return &User{persist: repository}
}

func (s *User) SignUp(ctx context.Context, usr *models.User) error {
	auth, ok := usr.Auth.(*au.BasicAuth)
	if !ok {
		return pkg.ErrInvalidTypeCast
	}

	existsUser, _ := s.findByEmail(ctx, auth)
	if existsUser != nil {
		return pkg.ErrEmailExists
	}

	err := auth.HashPassword()
	if err != nil {
		return fmt.Errorf("failed with hash password: %w", err)
	}

	return s.persist.Create(ctx, usr) //nolint:wrapcheck
}

func (s *User) SignIn(ctx context.Context, auth *au.BasicAuth) (*models.User, error) {
	foundUser, err := s.findByEmail(ctx, auth)
	if err != nil {
		return nil, pkg.ErrUserNotFound
		// return nil, fmt.Errorf("given email %s not found: %w", auth.Email, err)
	}

	err = foundUser.Auth.(*au.BasicAuth).ValidatePassword(auth.Password) //nolint:forcetypeassert
	if err != nil {
		return nil, pkg.ErrUserNotFound
		// return nil, fmt.Errorf("given password is not valid: %w", err)
	}

	return foundUser, nil
}

func (s *User) findByEmail(ctx context.Context, auth *au.BasicAuth) (*models.User, error) {
	auth.CleanCredentials() // не явная чистка значений
	foundUser, err := s.persist.FindByEmail(ctx, auth.Email)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return foundUser, nil
}
