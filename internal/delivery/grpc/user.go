package grpc

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	pbUser "github.com/UndeadDemidov/ya-pr-diplomb/gen_pb/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/delivery"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/telemetry"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pbUser.UserServiceServer = (*UserServer)(nil)

type UserServer struct {
	pbUser.UnimplementedUserServiceServer
	log        telemetry.AppLogger
	cfg        *config.App
	svc        delivery.User
	jwtManager auth.JWTManager
}

func NewUserServer(logger telemetry.AppLogger, config *config.App, service user.Service) *UserServer {
	return &UserServer{log: logger, cfg: config, svc: &service}
}

func (u *UserServer) SignUp(ctx context.Context, request *pbUser.SignUpRequest) (*emptypb.Empty, error) {
	usr := signinReq2User(request)
	l := u.log.With().Object("user", usr).Logger()
	l.Debug().Msg("signup user")

	if err := pkg.ValidateStruct(ctx, usr); err != nil {
		l.Err(err).Msg("pkg.ValidateStruct")
		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "ValidateStruct: %v", err)
	}

	err := u.svc.SignUp(ctx, usr)
	if err != nil {
		l.Err(err).Msg("UserServer.svc.SignUp")
		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "SignUp: %v", err)
	}

	l.Debug().Msg("user signed up successfully")
	return &emptypb.Empty{}, nil
}

func (u *UserServer) SignIn(ctx context.Context, request *pbUser.SignInRequest) (*pbUser.SignInResponse, error) {
	creds := credMsgToBasicAuth(request.GetCredentials())
	l := u.log.With().Object("creds", creds).Logger()
	l.Debug().Msg("signin user")

	if err := pkg.ValidateStruct(ctx, creds); err != nil {
		l.Err(err).Msg("pkg.ValidateStruct")

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "ValidateStruct: %v", err)
	}

	usr, err := u.svc.SignIn(ctx, creds)
	if err != nil {
		l.Err(err).Msg("UserServer.svc.SignIn")

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "SignIn: %v", err)
	}

	token, err := u.jwtManager.Generate(usr.UserUUID)
	if err != nil {
		l.Err(err).Stringer("user_uuid", usr.UserUUID).Msg("UserServer.jwtManager.Generate")

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "Generate: %v", err)
	}

	// TODO implement me
	// create session

	l.Debug().Stringer("user_uuid", usr.UserUUID).Msg("user signed on successfully")

	return &pbUser.SignInResponse{AccessToken: token, User: user2ProtoUser(usr)}, nil
}

func (u *UserServer) SignOut(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO implement me
	// invalidate session
	panic("implement me")
}

func (u *UserServer) mustEmbedUnimplementedUserServiceServer() { //nolint:unused
	// TODO implement me
	panic("implement me")
}
