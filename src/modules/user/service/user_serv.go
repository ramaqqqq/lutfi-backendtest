package service

import (
	"context"
	"folkatech-customerIdentity/src/modules/user/model"
	"folkatech-customerIdentity/src/pkg/utils"
)

type UserServ interface {
	CreateUser(ctx context.Context, request model.User) error
	UpdateUser(ctx context.Context, request model.User, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetList(ctx context.Context, filter model.FilterUser, pg *utils.PaginateQueryOffset) (response model.ListUsersResponse, err error)
	GetDetail(ctx context.Context, userID int64) (user model.UserResp, err error)
	GetByEmail(ctx context.Context, email string) (user model.User, err error)
}
