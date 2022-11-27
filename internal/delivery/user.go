package delivery

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
)

type User interface {
	SignIn(context.Context, *models.User) error
	SignOn(context.Context, *pkg.BasicAuth) (*models.User, error)
}
