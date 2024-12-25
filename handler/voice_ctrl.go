package handler

import (
	"context"
	"github.com/leor-w/kid"
	"role_ai/dto"
	"role_ai/infrastructure/web"
	"role_ai/services"
)

type VoiceController struct {
	srv *services.VoiceService `inject:""`
}

var voiceCtrl *VoiceController

func (ctrl *VoiceController) Provide(context.Context) any {
	voiceCtrl = ctrl
	return voiceCtrl
}

func GetVoiceList(ctx *kid.Context) any {
	var para dto.VoiceListReq
	if err := ctx.ShouldBindJSON(&para); err != nil {
		return web.ParamsErr(err)
	}
	resp, err := roleCtrl.srv.GetVoiceList(para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resp)
}
