package repo

import (
	"context"
	"folkatech-customerIdentity/src/modules/user/model"
	"folkatech-customerIdentity/src/pkg/utils"
)

type UserRepo interface {
	Save(ctx context.Context, user model.User) (*model.User, error)
	Update(ctx context.Context, userID int64, user model.User) error
	Delete(ctx context.Context, userID int64) error
	Find(ctx context.Context, filter model.FilterUser, pg *utils.PaginateQueryOffset) ([]model.User, int64, int64, error)
	FindByID(ctx context.Context, id int64) (model.User, error)
	FindByEmail(ctx context.Context, username string) (model.User, error)
}
