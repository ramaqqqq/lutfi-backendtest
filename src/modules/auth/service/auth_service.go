package service

import (
	"context"
	"folkatech-customerIdentity/src/modules/auth/model"
	user "folkatech-customerIdentity/src/modules/user/model"
)

type AuthService interface {
	Login(ctx context.Context, request model.AuthLogin) (model.AuthLoginResponse, error)
	Register(ctx context.Context, request model.AuthLogin) (user.UserResp, error)
}
