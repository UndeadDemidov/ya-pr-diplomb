package pkg

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info.
var validate *validator.Validate //nolint:gochecknoglobals

func init() {
	validate = validator.New()
}

// ValidateStruct validates struct fields.
func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s) //nolint:wrapcheck
}
