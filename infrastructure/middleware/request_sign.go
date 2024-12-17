package middleware

import (
	"context"
	errors2 "errors"
	"strconv"
	"time"

	"github.com/leor-w/kid"
	"github.com/leor-w/kid/errors"

	"role_ai/infrastructure/ecode"
	"role_ai/infrastructure/signer"
	"role_ai/infrastructure/utils"
	"role_ai/infrastructure/web"
	"role_ai/repos"
)

type RequestSign struct {
	//appRepo    repos.IAppRepository    `inject:""`
	signerRepo repos.ISignerRepository `inject:""`
}

var globalSign *RequestSign

func (sign *RequestSign) Provide(context.Context) any {
	globalSign = sign
	return globalSign
}

// VerifySign 验证签名
func VerifySign(ctx *kid.Context) {
	// 获取 appId 并检查值是否有效
	appId := ctx.GetHeader(signer.HeaderAppId.Value())
	if appId == "" {
		Resp(ctx, errors.New(ecode.AppIDRequired), nil)
		return
	}
	// 获取时间戳并检查值是否有效
	timestamp := ctx.GetHeader(signer.HeaderTimestamp.Value())
	if timestamp == "" {
		Resp(ctx, errors.New(ecode.TimestampRequired), nil)
		return
	}
	contentType := ctx.GetHeader(signer.HeaderContentType.Value())
	if contentType != signer.JSONContentType {
		Resp(ctx, errors.New(ecode.ContentTypeErr), nil)
		return
	}
	nonce := ctx.GetHeader(signer.HeaderNonce.Value())
	if !utils.CheckString(nonce) {
		Resp(ctx, errors.New(ecode.NonceRequired, errors2.New(nonce)), nil)
		return
	}
	signature := ctx.GetHeader(signer.HeaderSign.Value())
	if signature == "" {
		Resp(ctx, errors.New(ecode.SignatureRequired), nil)
		return
	}
	//// 获取应用信息
	//var app models.App
	//if err := globalSign.appRepo.GetOne(&finder.Finder{
	//	Model:     new(models.App),
	//	Wheres:    where.New().And(where.Eq("app_id", appId)),
	//	Recipient: &app,
	//	Debug:     true,
	//}); err != nil {
	//	ctx.JSON(200, web.Error(errors.New(ecode.AppNotFound)))
	//	ctx.Abort()
	//	return
	//}
	//// 验证签名
	//ok, err := signer.Verify(ctx.Request, app.Secret)
	//if err != nil || !ok {
	//	ctx.JSON(200, web.Error(errors.New(ecode.SignVerifyErr)))
	//	ctx.Abort()
	//	return
	//}
	//ctx.Set("app", &app) //记录应用标识
	ctx.Next()
}

func Resp(ctx *kid.Context, err error, data interface{}) {
	if err != nil {
		ctx.JSON(200, web.Error(err))
		ctx.Abort()
	} else {
		ctx.JSON(200, web.Success(data))
	}
}

// AntiReplayAttack 防重放攻击
func AntiReplayAttack(ctx *kid.Context) {
	// 检查请求是否超出时间限制
	timestamp, err := strconv.ParseInt(ctx.GetHeader(signer.HeaderTimestamp.Value()), 10, 64)
	if err != nil || time.Now().Unix()-timestamp > 300 {
		ctx.JSON(200, web.Error(errors.New(ecode.TimestampOutdated)))
		ctx.Abort()
		return
	}
	// 检查请求是否重放
	nonce := ctx.GetHeader(signer.HeaderNonce.Value())
	if err := globalSign.signerRepo.CheckNonce(nonce); err != nil {
		ctx.JSON(200, web.Error(errors.New(ecode.NonceReplayAttack)))
		ctx.Abort()
		return
	}
	// 保存 nonce
	if err := globalSign.signerRepo.SaveNonce(nonce); err != nil {
		ctx.JSON(200, web.Error(errors.New(ecode.NonceSaveErr)))
		ctx.Abort()
		return
	}
	ctx.Next()
}

// AES-256-CBC 加密

// Idempotence 防幂等攻击
func Idempotence(ctx *kid.Context) {
	// 检查请求是否重复
	idempotent, err := utils.GenerateFingerprint(ctx.Request)
	if err != nil {
		ctx.JSON(200, web.Error(errors.New(ecode.FingerprintGenerateErr)))
		ctx.Abort()
		return
	}

	if err := globalSign.signerRepo.CheckIdempotent(idempotent); err != nil {
		ctx.JSON(200, web.Error(errors.New(ecode.IdempotentReplayAttack)))
		ctx.Abort()
		return
	}

	// 保存 idempotent
	if err := globalSign.signerRepo.SaveIdempotent(idempotent); err != nil {
		ctx.JSON(200, web.Error(errors.New(ecode.IdempotentSaveErr)))
		ctx.Abort()
		return
	}
	ctx.Next()

	// 最后删除 idempotent
	_ = globalSign.signerRepo.DeleteIdempotent(idempotent)
}
