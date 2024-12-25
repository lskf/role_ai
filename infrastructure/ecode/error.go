package ecode

import (
	"github.com/leor-w/kid/errors"
)

var (
	OK = &errors.Status{Code: 200, Message: "成功"}

	// 100xx 基本常用错误
	ReqParamMissErr        = &errors.Status{Code: 10001, Message: "必填参数缺失或参数及类型错误"}
	ReqParamInvalidErr     = &errors.Status{Code: 10002, Message: "参数错误"}
	InternalErr            = &errors.Status{Code: 10003, Message: "内部错误"}
	NotFoundRouteErr       = &errors.Status{Code: 10005, Message: "无效的请求链接"}
	DatabaseErr            = &errors.Status{Code: 10009, Message: "数据操作错误"}
	DataProcessingErr      = &errors.Status{Code: 10015, Message: "数据处理错误"}
	NotFoundOperatorErr    = &errors.Status{Code: 10017, Message: "未找到用户信息"}
	DataNotExist           = &errors.Status{Code: 10018, Message: "数据不存在"}
	SMSSendErr             = &errors.Status{Code: 10019, Message: "发送短信验证码错误"}
	SendSMSCodeIntervalErr = &errors.Status{Code: 10020, Message: "您的操作过于频繁，请稍后再试。"}
	LoginAuthCodeErr       = &errors.Status{Code: 10021, Message: "验证码错误"}

	NotFoundTokenErr      = &errors.Status{Code: 20003, Message: "登录令牌缺失"}
	AuthFailedErr         = &errors.Status{Code: 20004, Message: "验证失败"}
	UserNotFoundErr       = &errors.Status{Code: 20005, Message: "用户不存在"}
	TokenNotExistErr      = &errors.Status{Code: 20007, Message: "登录令牌已失效，请重新登录"}
	GenerateSecretErr     = &errors.Status{Code: 20011, Message: "生成密钥失败"}
	GenerateEncryptKeyErr = &errors.Status{Code: 20012, Message: "生成加密密钥失败"}

	// 300xx 签名验证错误
	AppIDRequired          = &errors.Status{Code: 30001, Message: "AppID 不能为空"}
	AppNotFound            = &errors.Status{Code: 30002, Message: "AppID 不存在"}
	SignVerifyErr          = &errors.Status{Code: 30003, Message: "签名验证失败"}
	TimestampOutdated      = &errors.Status{Code: 30004, Message: "请求已过期"}
	NonceReplayAttack      = &errors.Status{Code: 30005, Message: "nonce 已存在，疑似重放攻击"}
	NonceSaveErr           = &errors.Status{Code: 30006, Message: "请求处理失败，请稍后再试"}
	FingerprintGenerateErr = &errors.Status{Code: 30007, Message: "请求处理失败，请稍后再试"}
	IdempotentReplayAttack = &errors.Status{Code: 30008, Message: "请求已处理，请勿重复提交"}
	IdempotentSaveErr      = &errors.Status{Code: 30009, Message: "请求处理失败，请稍后再试"}
	TimestampRequired      = &errors.Status{Code: 30010, Message: "Timestamp 不能为空"}
	ContentTypeErr         = &errors.Status{Code: 30011, Message: "Content-Type 错误"}
	NonceRequired          = &errors.Status{Code: 30012, Message: "Nonce 必须为字母与数字，长度不能小于 8 且不能超过 16"}
	SignatureRequired      = &errors.Status{Code: 30013, Message: "X-Tsign-Sign 不能为空"}

	// 400xx 登录错误
	UserNotFound    = &errors.Status{Code: 40001, Message: "用户不存在"}
	AdminNotFound   = &errors.Status{Code: 40002, Message: "管理员不存在"}
	PasswordErr     = &errors.Status{Code: 40003, Message: "用户不存在或密码错误"}
	TokenErr        = &errors.Status{Code: 40004, Message: "令牌生成失败"}
	GenCaptchaErr   = &errors.Status{Code: 40005, Message: "验证码生成失败"}
	CaptchaErr      = &errors.Status{Code: 40006, Message: "验证码错误"}
	LogoutErr       = &errors.Status{Code: 40007, Message: "退出登录失败"}
	TokenIsEmptyErr = &errors.Status{Code: 40008, Message: "令牌为空"}

	//500xx 角色错误
	RoleNotExistErr = &errors.Status{Code: 50001, Message: "角色不存在"}
	VoiceNotFound   = &errors.Status{Code: 50002, Message: "声音不存在"}
	PublicChangeErr = &errors.Status{Code: 50003, Message: "公开角色不能改为私密"}
)
