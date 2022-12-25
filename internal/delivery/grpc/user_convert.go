package grpc

import (
	"github.com/UndeadDemidov/ya-pr-diplomb/gen_pb"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
)

func credMsgToBasicAuth(c *gen_pb.Credentials) *auth.BasicAuth {
	return &auth.BasicAuth{
		Email:    c.GetEmail(),
		Password: c.GetPassword(),
	}
}

func signinReq2User(r *gen_pb.SignUpRequest) *models.User {
	return models.NewUser(credMsgToBasicAuth(r.GetCredentials()))
}

func user2ProtoUser(usr *models.User) *gen_pb.User {
	return &gen_pb.User{Uuid: usr.UserUUID.String()}
}

// func signonMsgToModel(r *userService.SignOnRequest) *models.User {
// 	return &models.User{
// 		Auth: credMsgToBasicAuth(r.Credentials),
// 	}
// }
