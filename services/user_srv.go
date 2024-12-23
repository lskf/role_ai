package services

import (
	"context"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/updater"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"role_ai/common"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/models"
	"role_ai/repos"
	"time"
)

type UserService struct {
	userRepo repos.IUserRepository `inject:""`
}

func (srv *UserService) Provide(_ context.Context) interface{} {
	return srv
}

func (srv *UserService) GetUserList() {}

func (srv *UserService) GetUserDetail(uid int64) (*dto.User, error) {
	user := models.User{}
	err := srv.userRepo.GetOne(&finder.Finder{
		Model:     new(models.User),
		Wheres:    where.New().And(where.Eq("uid", uid)),
		Recipient: &user,
	})
	if err != nil {
		return nil, errors.New(ecode.UserNotFound)
	}
	userDetail := dto.User{}
	err = models.Copy(&userDetail, &user)
	if err != nil {
		return nil, errors.New(ecode.DataProcessingErr)
	}
	userDetail.BirthdayStr = userDetail.Birthday.Format(common.TimeFormatToDateTime)
	return &userDetail, nil
}

func (srv *UserService) CreateUser() {}

func (srv *UserService) UpdateUser(uid int64, data *dto.User) (err error) {
	user := models.User{}
	err = srv.userRepo.GetOne(&finder.Finder{
		Model:     new(models.User),
		Wheres:    where.New().And(where.Eq("uid", uid)),
		Recipient: &user,
	})
	if err != nil {
		return errors.New(ecode.UserNotFound, err)
	}
	err = models.Copy(&user, &data)
	if err != nil {
		return errors.New(ecode.DataProcessingErr, err)
	}
	err = srv.userRepo.Update(&updater.Updater{
		Model:  new(models.User),
		Wheres: where.New().And(where.Eq("uid", uid)),
		Fields: map[string]interface{}{
			"account":    data.Account,
			"nick_name":  data.NickName,
			"avatar":     data.Avatar,
			"gender":     data.Gender,
			"birthday":   data.Birthday,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return errors.New(ecode.DatabaseErr, err)
	}
	return nil
}

func (srv *UserService) DeleteUser() {}
