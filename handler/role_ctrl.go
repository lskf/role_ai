package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/leor-w/kid"
	"role_ai/dto"
	"role_ai/infrastructure/web"
	"role_ai/models"
	"role_ai/services"
	"strings"
)

type RoleController struct {
	srv *services.RoleService `inject:""`
}

var roleCtrl *RoleController

func (ctrl *RoleController) Provide(context.Context) any {
	roleCtrl = ctrl
	return roleCtrl
}

// GetRoleList
// @Description: 获取角色列表
// @param ctx
// @return any
func GetRoleList(ctx *kid.Context) any {
	var para dto.RoleListReq
	if err := ctx.ShouldBindQuery(&para); err != nil {
		return web.ParamsErr(err)
	}
	resp, err := roleCtrl.srv.GetRoleList(para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resp)
}

// GetRoleDetail
// @Description: 获取角色详情
// @param ctx
// @return any
func GetRoleDetail(ctx *kid.Context) any {
	if err := ctx.ValidFiled(kid.Rules{
		"id": "required",
	}); err != nil {
		return web.ParamsErr(err)
	}
	id := ctx.FindInt64("id")

	res, err := roleCtrl.srv.GetDetailById(id)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(res)
}

// CreateRole
// @Description: 创建角色
// @param ctx
// @return any
func CreateRole(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var req dto.CreateRoleReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	if len(req.TagArray) > 0 {
		req.Tag = strings.Join(req.TagArray, ";")
	}
	gamificationStr, err := json.Marshal(req.GamificationObj)
	if err != nil {
		return web.ParamsErr(err)
	}
	req.Gamification = string(gamificationStr)
	id, err := roleCtrl.srv.CreateRole(&user, req)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(dto.IdField{Id: id})
}

// UpdateRole
// @Description: 编辑角色
// @param ctx
// @return any
func UpdateRole(ctx *kid.Context) any {
	if err := ctx.ValidFiled(kid.Rules{
		"id": "required",
	}); err != nil {
		return web.ParamsErr(err)
	}
	id := ctx.FindInt64("id")

	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var req dto.UpdateRoleResp
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	if len(req.TagArray) > 0 {
		req.Tag = strings.Join(req.TagArray, ";")
	}
	gamificationStr, err := json.Marshal(req.GamificationObj)
	if err != nil {
		return web.ParamsErr(err)
	}
	req.Gamification = string(gamificationStr)
	req.Id = id
	err = roleCtrl.srv.UpdateRole(user, req)
	if err != nil {
		return web.Error(err)
	}
	return web.Success()
}

func GetRoleAvatarSetting(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var req dto.AiCreateRoleReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	resq, err := roleCtrl.srv.GetRoleAvatarSetting(&req)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resq)
}

func GetRoleSetting(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var req dto.AiCreateRoleReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	resq, err := roleCtrl.srv.GetRoleSetting(&req)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resq)
}

func CreateRoleAvatar(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var req dto.CreateRoleAvatarReq
	if err := ctx.Valid(&req); err != nil {
		return web.ParamsErr(err)
	}
	resp, err := roleCtrl.srv.CreateRoleAvatar(user.Uid, &req)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resp)
}

func GetRoleAvatarHistory(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	promptId := ctx.FindString("prompt_id")
	if promptId == "" {
		return web.ParamsErr(errors.New("prompt_id 为空"))
	}
	resp, err := roleCtrl.srv.GetRoleAvatarHistory(promptId)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(resp)
}

func GetRoleAvatar(ctx *kid.Context) any {
	user := models.User{}
	err := ctx.GetBindUser(&user)
	if err != nil {
		return web.Unauthorized(err)
	}
	var req dto.GetViewReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return web.ParamsErr(err)
	}
	resp, err := roleCtrl.srv.GetRoleAvatar(&req)
	if err != nil {
		return web.Error(err)
	}
	return resp
}
