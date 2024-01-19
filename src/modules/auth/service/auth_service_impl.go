package service

import (
	"context"
	"errors"
	"folkatech-customerIdentity/src/middleware"
	"folkatech-customerIdentity/src/modules/auth/model"
	user "folkatech-customerIdentity/src/modules/user/model"
	u "folkatech-customerIdentity/src/modules/user/service"
	"folkatech-customerIdentity/src/pkg/helpers"
	"strings"

	"github.com/go-playground/validator"
)

type AuthServiceImpl struct {
	UserService u.UserServ
	Validate    *validator.Validate
}

func NewAuthService(userService u.UserServ, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserService: userService,
		Validate:    validate,
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, request model.AuthLogin) (model.AuthLoginResponse, error) {
	if s.UserService == nil {
		return model.AuthLoginResponse{}, errors.New("UserService is not initialized")
	}

	user, err := s.UserService.GetByEmail(ctx, strings.ToLower(request.EmailAddress))
	if err != nil {
		return model.AuthLoginResponse{}, err
	}

	isPasswordValid := helpers.CheckPasswordHash(request.Password, user.Password)
	if !isPasswordValid {
		err = errors.New("wrong email or password")
		return model.AuthLoginResponse{}, err
	}

	tokens, _ := middleware.CreateToken(string(user.ID), user.Username, user.EmailAddress)

	loginResponse := model.AuthLoginResponse{
		User: model.AuthLogin{
			ID:             user.ID,
			Username:       user.Username,
			EmailAddress:   user.EmailAddress,
			AccountNumber:  user.AccountNumber,
			IdentityNumber: user.IdentityNumber,
		},
		Tokens: model.JwtModel{
			AccessToken:  tokens["accessToken"],
			RefreshToken: tokens["refreshToken"],
		},
	}

	return loginResponse, nil
}

func (s *AuthServiceImpl) Register(ctx context.Context, request model.AuthLogin) (user.UserResp, error) {
	err := s.UserService.CreateUser(ctx, user.User(request))

	if err != nil {
		return user.UserResp{}, err
	}

	return user.UserResp{
		Username:       request.Username,
		EmailAddress:   request.EmailAddress,
		AccountNumber:  request.AccountNumber,
		IdentityNumber: request.IdentityNumber,
	}, nil
}
