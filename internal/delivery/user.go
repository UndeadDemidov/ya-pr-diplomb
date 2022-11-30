package delivery

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
)

type User interface {
	SignUp(context.Context, *models.User) error
	SignIn(context.Context, *auth.BasicAuth) (*models.User, error)
}
