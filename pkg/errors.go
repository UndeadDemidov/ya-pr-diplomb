package pkg

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

// Service errors.
var (
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidTypeCast = errors.New("can't cast interface to given type")
)

// ParseGRPCErrStatusCode parses error and get GRPC code.
func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	// case errors.Is(err, redis.Nil):
	// 	return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrEmailExists):
		return codes.AlreadyExists
	// case errors.Is(err, ErrNoCtxMetaData):
	// 	return codes.Unauthenticated
	// case errors.Is(err, ErrInvalidSessionId):
	// 	return codes.PermissionDenied
	case strings.Contains(err.Error(), "Validate"):
		return codes.InvalidArgument
	case strings.Contains(err.Error(), "redis"):
		return codes.NotFound
	}

	return codes.Internal
}
