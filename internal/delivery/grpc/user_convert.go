package grpc

import (
	userService "github.com/UndeadDemidov/ya-pr-diplomb/gen_pb/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/auth"
)

func credMsgToBasicAuth(c *userService.Credentials) *auth.BasicAuth {
	return &auth.BasicAuth{
		Email:    c.GetEmail(),
		Password: c.GetPassword(),
	}
}

func signinMsgToUser(r *userService.SignInRequest) *models.User {
	return models.NewUser(credMsgToBasicAuth(r.User.GetCredentials()))
}

// func signonMsgToModel(r *userService.SignOnRequest) *models.User {
// 	return &models.User{
// 		Auth: credMsgToBasicAuth(r.Credentials),
// 	}
// }
