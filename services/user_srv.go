package services

import (
	"context"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/models"
	"role_ai/repos"
)

type UserService struct {
	userRepo repos.IUserRepository `inject:""`
}

func (srv *UserService) Provide(_ context.Context) interface{} {
	return srv
}

func (srv *UserService) GetUserList() {}

func (srv *UserService) GetUserDetail(uid int64) (userDetail *dto.User, err error) {
	user := models.User{}
	err = srv.userRepo.GetOne(&finder.Finder{
		Model:     new(models.User),
		Wheres:    where.New().And(where.Eq("uid", uid)),
		Recipient: &user,
	})
	if err != nil {
		return nil, errors.New(ecode.UserNotFound)
	}
	err = models.Copy(&user, &userDetail)
	return nil, nil
}

func (srv *UserService) CreateUser() {}

func (srv *UserService) UpdateUser() {}

func (srv *UserService) DeleteUser() {}
