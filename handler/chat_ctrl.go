package handler

import (
	"context"
	"github.com/leor-w/kid"
	"role_ai/dto"
	"role_ai/infrastructure/web"
	"role_ai/models"
	"role_ai/services"
)

type ChatController struct {
	srv *services.ChatService `inject:""`
}

var chatCtrl *ChatController

func (ctrl *ChatController) Provide(context.Context) any {
	chatCtrl = ctrl
	return chatCtrl
}

func Chat(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var para dto.ChatReq
	if err := ctx.Valid(&para); err != nil {
		return web.ParamsErr(err)
	}
	resp, err := chatCtrl.srv.Chat(user.Uid, para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resp)
}

func ChatList(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var para dto.ChatListReq
	if err = ctx.ShouldBindQuery(&para); err != nil {
		return web.ParamsErr(err)
	}
	res, err := chatCtrl.srv.GetList(user.Uid, para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(res)
}

func ChatDetail(ctx *kid.Context) any {
	return nil
}

func EditAnswer(ctx *kid.Context) any {
	return nil
}

func DelChat(ctx *kid.Context) any {
	return nil
}

func GetChatHistoryList(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var para dto.ChatHistoryListReq
	if err = ctx.ShouldBindQuery(&para); err != nil {
		return web.ParamsErr(err)
	}
	res, err := chatCtrl.srv.GetHistoryList(user.Uid, para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(res)
}
