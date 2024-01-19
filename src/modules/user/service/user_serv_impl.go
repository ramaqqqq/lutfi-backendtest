package service

import (
	"context"
	"folkatech-customerIdentity/src/modules/user/model"
	"folkatech-customerIdentity/src/modules/user/repo"
	"folkatech-customerIdentity/src/pkg/helpers"
	"folkatech-customerIdentity/src/pkg/utils"
)

type UserServImpl struct {
	Repository repo.UserRepo
}

func NewUserService(userRepository repo.UserRepo) UserServ {
	return &UserServImpl{
		Repository: userRepository,
	}
}

func (s *UserServImpl) CreateUser(ctx context.Context, request model.User) error {
	hashed, err := helpers.Hash(request.Password)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	_, err = s.Repository.Save(ctx, model.User{
		Username:       request.Username,
		AccountNumber:  request.AccountNumber,
		EmailAddress:   request.EmailAddress,
		IdentityNumber: request.IdentityNumber,
		Password:       string(hashed),
	})

	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	return nil
}

func (s *UserServImpl) UpdateUser(ctx context.Context, request model.User, id int64) error {
	user, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	user.Username = request.Username
	user.AccountNumber = request.AccountNumber
	user.EmailAddress = request.EmailAddress
	user.IdentityNumber = request.IdentityNumber

	if err = s.Repository.Update(ctx, id, model.User{
		Username:       user.Username,
		EmailAddress:   user.EmailAddress,
		AccountNumber:  user.AccountNumber,
		IdentityNumber: user.IdentityNumber,
	}); err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	return nil
}

func (s *UserServImpl) DeleteUser(ctx context.Context, id int64) error {
	if _, err := s.Repository.FindByID(ctx, id); err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	if err := s.Repository.Delete(ctx, id); err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	return nil
}

func (s *UserServImpl) GetList(ctx context.Context, filter model.FilterUser, pg *utils.PaginateQueryOffset) (response model.ListUsersResponse, err error) {
	data, totalData, totalPage, err := s.Repository.Find(ctx, filter, pg)
	if err != nil {
		return
	}

	if totalData == 0 {
		response.TotalPage = 0
		response.TotalItem = 0
		response.User = []model.UserResp{}
		return
	}

	var userRespData []model.UserResp
	for _, user := range data {
		userResp := model.UserResp{
			Username:       user.Username,
			AccountNumber:  user.AccountNumber,
			EmailAddress:   user.EmailAddress,
			IdentityNumber: user.IdentityNumber,
			BaseTime:       user.BaseTime,
		}
		userRespData = append(userRespData, userResp)
	}

	response.User = userRespData
	response.TotalPage = totalPage
	response.TotalItem = totalData

	return
}

func (s *UserServImpl) GetDetail(ctx context.Context, userID int64) (user model.UserResp, err error) {
	userModel, err := s.Repository.FindByID(ctx, userID)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return
	}

	user = model.UserResp{
		Username:       userModel.Username,
		AccountNumber:  userModel.AccountNumber,
		EmailAddress:   userModel.EmailAddress,
		IdentityNumber: userModel.IdentityNumber,
	}

	return
}

func (s *UserServImpl) GetByEmail(ctx context.Context, email string) (user model.User, err error) {
	if user, err = s.Repository.FindByEmail(ctx, email); err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return
	}

	return
}
