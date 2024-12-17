package services

import (
	"context"
	"github.com/leor-w/kid/database/repos/creator"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"github.com/leor-w/kid/guard"
	"github.com/leor-w/kid/plugin/captcha"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/models"
	"role_ai/repos"
	"strings"
	"time"
)

type LoginService struct {
	userRepo repos.IUserRepository `inject:""`

	captcha *captcha.Captcha `inject:""`
	guard   guard.Guard      `inject:""`
}

func (srv *LoginService) Provide(_ context.Context) interface{} {
	return srv
}

// GetLoginCaptcha 获取图形验证码
func (srv *LoginService) GetLoginCaptcha() (*dto.CaptchaResp, error) {
	key, code, err := srv.captcha.Generate()
	if err != nil {
		return nil, errors.New(ecode.GenCaptchaErr, err)
	}
	return &dto.CaptchaResp{
		Key:     key,
		Captcha: code,
	}, nil
}

func (srv *LoginService) Login(phone string) (*models.User, error) {
	//获取用户信息
	var user models.User
	_ = srv.userRepo.GetOne(&finder.Finder{
		Model:     new(models.User),
		Wheres:    where.New().And(where.Eq("phone", phone)),
		Recipient: &user,
	})
	//判断用户是否存在
	if user.Id <= 0 || user.Uid <= 0 {
		//注册
		user.Phone = phone
		err := srv.Registered(&user)
		if err != nil {
			return nil, err
		}
	}
	token, err := srv.guard.License(&guard.User{
		Uid:  user.Uid,
		Type: guard.General,
	})
	if err != nil {
		return nil, errors.New(ecode.TokenErr, err)
	}
	user.Token = token
	return &user, nil
}

// Registered
// @Description: 注册
// @receiver srv
// @param phone
// @return *models.User
// @return error
func (srv *LoginService) Registered(user *models.User) error {
	//获取唯一Uid
	user.Uid = srv.userRepo.GetUniqueID(&finder.Finder{
		Model:  new(models.User),
		Wheres: where.New().And(where.Eq("uid", nil)),
	}, 10000000, 99999999, 0, 0)
	user.Account = user.Phone
	user.NickName = user.Phone
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Menber = models.UserMenberNormal
	err := srv.userRepo.Create(&creator.Creator{
		Data: &user,
	})
	if err != nil {
		return errors.New(ecode.InternalErr, err)
	}
	return nil
}

// LoginByPassword 账号密码登录
func (srv *LoginService) LoginByPassword(phone, password, key, captcha string) (*models.User, error) {
	// 校验图形验证码
	if !srv.captcha.Verify(key, captcha, true) {
		return nil, errors.New(ecode.CaptchaErr)
	}

	// 校验账号密码
	var user models.User
	if err := srv.userRepo.GetOne(&finder.Finder{
		Model:     new(models.User),
		Wheres:    where.New().And(where.Eq("phone", phone)),
		Recipient: &user,
	}); err != nil {
		return nil, errors.New(ecode.UserNotFoundErr)
	}
	//if !user.CheckPassword(password) {
	//	return nil, errors.New(ecode.PasswordErr)
	//}
	token, err := srv.guard.License(&guard.User{
		Uid:  user.Uid,
		Type: guard.General,
	})
	if err != nil {
		return nil, errors.New(ecode.TokenErr)
	}
	user.Token = token
	return &user, nil
}

// Logout 退出登录
func (srv *LoginService) Logout(token string, uid int64) error {
	if len(strings.Trim(token, " /n/t")) > 0 {
		if err := srv.guard.Cancellation(token); err != nil {
			return errors.New(ecode.LogoutErr, err)
		}
		return nil
	}
	return errors.New(ecode.TokenIsEmptyErr)
}
