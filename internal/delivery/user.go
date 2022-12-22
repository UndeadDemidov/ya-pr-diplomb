package delivery

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	_ "github.com/golang/mock/mockgen/model" //
)

//go:generate mockgen -destination=./mocks/mock_user.go . User

type User interface {
	SignUp(context.Context, *models.User) error
	SignIn(context.Context, *auth.BasicAuth) (*models.User, error)
}
