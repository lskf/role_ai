package handler

import (
	"context"

	"github.com/leor-w/kid"

	"role_ai/dto"
	"role_ai/infrastructure/web"
	"role_ai/models"
	"role_ai/services"
)

type LoginController struct {
	srv    *services.LoginService `inject:""`
	smsSrv *services.SmsService   `inject:""`
}

var loginCtrl *LoginController

func (ctrl *LoginController) Provide(context.Context) any {
	loginCtrl = ctrl
	return loginCtrl
}

// GetLoginCaptcha
// @Description: 获取图形验证码
// @param *kid.Context
// @return any
func GetLoginCaptcha(*kid.Context) any {
	resp, err := loginCtrl.srv.GetLoginCaptcha()
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resp)
}

// SendSmsCode
// @Description: 发送短信验证码
// @param ctx
// @return any
func SendSmsCode(ctx *kid.Context) any {
	var req dto.SendSmsReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	err := loginCtrl.smsSrv.SendCode(req.Phone)
	if err != nil {
		return web.Error(err)
	}
	return web.Success()
}

// LoginBySms
// @Description: 短信验证码登录
// @param ctx
// @return any
func LoginBySms(ctx *kid.Context) any {
	var req dto.LoginBySmsReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	//校验短信验证码
	ok, err := loginCtrl.smsSrv.VerifyCode(req.Phone, req.Code)
	if !ok || err != nil {
		return web.Error(err)
	}
	//登录
	user, err := loginCtrl.srv.Login(req.Phone)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(dto.BuildLoginResp(user))
}

// LoginByPassword
// @Description: 密码登录
// @param ctx
// @return any
func LoginByPassword(ctx *kid.Context) any {
	var req dto.LoginByPasswordReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	user, err := loginCtrl.srv.LoginByPassword(req.Phone, req.Password, req.Key, req.Captcha)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(dto.BuildLoginResp(user))
}

// Logout
// @Description: 退出登录
// @param ctx
// @return any
func Logout(ctx *kid.Context) any {
	var admin models.Admin
	if err := ctx.GetBindUser(&admin); err != nil {
		return web.Unauthorized(err)
	}
	if err := loginCtrl.srv.Logout(admin.Token, 0); err != nil {
		return web.Error(err)
	}
	return web.Success()
}
