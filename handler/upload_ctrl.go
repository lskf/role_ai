package handler

import (
	"context"
	"github.com/leor-w/kid"
	"role_ai/dto"
	"role_ai/infrastructure/web"
	"role_ai/models"
	"role_ai/services"
)

type UploadController struct {
	srv *services.UploadService `inject:""`
}

var uploadCtrl *UploadController

func (ctrl *UploadController) Provide(context.Context) any {
	uploadCtrl = ctrl
	return uploadCtrl
}

func Upload(ctx *kid.Context) any {
	//判断是否登录
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	path, err := uploadCtrl.srv.UploadFormFile(ctx, user.Uid)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(dto.UploadFileResp{FilePath: path})
}

func UploadByUrl(ctx *kid.Context) any {
	//判断是否登录
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var req dto.UploadFileByUrlReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	path, err := uploadCtrl.srv.UploadUrlFile(user.Uid, req.Url, req.FilePathType)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(dto.UploadFileResp{FilePath: path})
}
