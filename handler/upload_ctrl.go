package handler

import (
	"context"
	"github.com/leor-w/kid"
	"role_ai/infrastructure/web"
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

func GetUploadToken(ctx *kid.Context) interface{} {
	return web.Success(uploadCtrl.srv.GetQiniuUploadToken())
}
