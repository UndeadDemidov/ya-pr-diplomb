package interceptors

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a server interceptor for authentication and authorization.
type AuthInterceptor struct {
	log        *zerolog.Logger
	jwtManager *auth.JWTManager
}

// NewAuthInterceptor returns a new auth interceptor.
func NewAuthInterceptor(jwtManager *auth.JWTManager, logger *zerolog.Logger) *AuthInterceptor {
	return &AuthInterceptor{logger, jwtManager}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC.
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		i.log.Info().Msgf("--> unary i: %s", info.FullMethod)

		err := i.authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC.
func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		i.log.Info().Msgf("--> stream i: %s", info.FullMethod)

		err := i.authorize(stream.Context())
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := i.jwtManager.Verify(accessToken)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	// ToDo claims put in context
	i.log.Info().Msgf("access token claims %v", claims)

	return status.Error(codes.PermissionDenied, "no permission to access this RPC") //nolint:wrapcheck
}
