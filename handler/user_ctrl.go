package handler

import (
	"context"
	"github.com/leor-w/kid"
	"github.com/leor-w/kid/errors"
	"role_ai/dto"
	"role_ai/infrastructure/ecode"
	"role_ai/infrastructure/web"
	"role_ai/models"
	"role_ai/services"
	"time"
)

type UserController struct {
	srv *services.UserService `inject:""`
}

var userCtrl *UserController

func (ctrl *UserController) Provide(context.Context) any {
	userCtrl = ctrl
	return userCtrl
}

func GetUserDetail(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	res, err := userCtrl.srv.GetUserDetail(user.Uid)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(res)
}

func UpdateUser(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	var req dto.User
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	req.Birthday, err = time.Parse(dto.TimeFormatToDateTime, req.BirthdayStr)
	if err != nil {
		return web.ParamsErr(errors.New(ecode.ReqParamInvalidErr, err))
	}
	err = userCtrl.srv.UpdateUser(user.Uid, &req)
	if err != nil {
		return web.Error(err)
	}
	return web.Success()
}
