package dto

import "role_ai/models"

// CaptchaResp 验证码响应
type CaptchaResp struct {
	Key     string `json:"key"`
	Captcha string `json:"captcha"`
}

type SendSmsReq struct {
	Phone string `json:"phone" validate:"required"` // 手机号
}

// LoginBySmsReq 短信验证码登录
type LoginBySmsReq struct {
	Phone string `json:"phone" validate:"required"` // 手机号
	Code  string `json:"code" validate:"required"`  //短信验证码
}

// LoginByPasswordReq 账号密码登录请求
type LoginByPasswordReq struct {
	Phone    string `json:"phone" validate:"required"`    // 手机号
	Password string `json:"password" validate:"required"` // 密码
	Key      string `json:"key" validate:"required"`      // 验证码key
	Captcha  string `json:"captcha" validate:"required"`  // 验证码
}

// LoginResp 登录响应
type LoginResp struct {
	Uid     int64  `json:"uid"`
	Phone   string `json:"phone"`
	Account string `json:"account"`
	Avatar  string `json:"avatar"`
	Token   string `json:"token"`
}

func BuildLoginResp(user *models.User) *LoginResp {
	return &LoginResp{
		Uid:     user.Uid,
		Phone:   user.Phone,
		Account: user.Account,
		Avatar:  user.Avatar,
		Token:   user.Token,
	}

}
