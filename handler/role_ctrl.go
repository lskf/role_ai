package handler

import (
	"context"
	"encoding/json"
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
	if err := ctx.ShouldBindJSON(&para); err != nil {
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
