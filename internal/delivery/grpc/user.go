package grpc

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	pbUser "github.com/UndeadDemidov/ya-pr-diplomb/gen_pb/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/delivery"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pbUser.UserServiceServer = (*UserServer)(nil)

type UserServer struct {
	pbUser.UnimplementedUserServiceServer
	log        pkg.Logger
	cfg        *config.App
	svc        delivery.User
	jwtManager pkg.JWTManager
}

func NewUserServer(logger pkg.Logger, config *config.App, service user.Service) *UserServer {
	return &UserServer{log: logger, cfg: config, svc: &service}
}

func (u *UserServer) SignIn(ctx context.Context, request *pbUser.SignInRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc.UserServer.SignIn")
	defer span.Finish()

	usr := signinMsgToUser(request)
	if err := pkg.ValidateStruct(ctx, usr); err != nil {
		u.log.Errorf("ValidateStruct: %v", err)

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "ValidateStruct: %v", err)
	}

	err := u.svc.SignIn(ctx, usr)
	if err != nil {
		u.log.Errorf("service.SignIn: %v", err)

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "SignIn: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (u *UserServer) SignOn(ctx context.Context, request *pbUser.SignOnRequest) (*pbUser.SignOnResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc.UserServer.SignOn")
	defer span.Finish()

	creds := credMsgToBasicAuth(request.GetCredentials())
	if err := pkg.ValidateStruct(ctx, creds); err != nil {
		u.log.Errorf("ValidateStruct: %v", err)

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "ValidateStruct: %v", err)
	}

	usr, err := u.svc.SignOn(ctx, creds)
	if err != nil {
		u.log.Errorf("service.SignOn: %v", err)

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "SignOn: %v", err)
	}

	token, err := u.jwtManager.Generate(usr.UserUUID)
	if err != nil {
		u.log.Errorf("jwtManager.Generate: %v", err)

		return nil, status.Errorf(pkg.ParseGRPCErrStatusCode(err), "Generate: %v", err)
	}

	// TODO implement me
	// create session

	return &pbUser.SignOnResponse{AccessToken: token}, nil
}

func (u *UserServer) SignOut(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	// span, ctx := opentracing.StartSpanFromContext(ctx, "grpc.UserServer.SignOut")
	// defer span.Finish()

	// TODO implement me
	// invalidate session
	panic("implement me")
}

// func (u UserServer) mustEmbedUnimplementedUserServiceServer() {
// 	// TODO implement me
// 	panic("implement me")
// }
