package handler

import (
	"context"
	"github.com/leor-w/kid"
	"role_ai/dto"
	"role_ai/infrastructure/web"
	"role_ai/models"
	"role_ai/services"
)

type AdminController struct {
	srv *services.AdminService `inject:""`
}

var adminCtrl *AdminController

func (ctrl *AdminController) Provide(context.Context) any {
	adminCtrl = ctrl
	return adminCtrl
}

func GetAdminList(ctx *kid.Context) any {
	para := dto.AdminListReq{}
	if err := ctx.ShouldBindQuery(&para); err != nil {
		return web.ParamsErr(err)
	}

	res, err := adminCtrl.srv.GetList(&para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(res)
}

func GetAdminDetail(ctx *kid.Context) any {
	if err := ctx.ValidFiled(kid.Rules{
		"uid": "required",
	}); err != nil {
		return web.ParamsErr(err)
	}
	uid := ctx.FindInt64("uid")

	res, err := adminCtrl.srv.Detail(uid)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(res)
}

func CreateAdmin(ctx *kid.Context) any {
	var para dto.AdminCreateReq
	if err := ctx.Valid(&para); err != nil {
		return web.ParamsErr(err)
	}

	var admin models.Admin
	if err := ctx.GetBindUser(&admin); err != nil {
		return web.Unauthorized(err)
	}
	uid, err := adminCtrl.srv.Create(&admin, &para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(uid)
}

func UpdateAdmin(ctx *kid.Context) any {
	if err := ctx.ValidFiled(kid.Rules{
		"uid": "required",
	}); err != nil {
		return web.ParamsErr(err)
	}
	uid := ctx.FindInt64("uid")
	var para dto.AdminUpdateReq
	if err := ctx.Valid(&para); err != nil {
		return web.ParamsErr(err)
	}
	para.Uid = uid

	var admin models.Admin
	if err := ctx.GetBindUser(&admin); err != nil {
		return web.Unauthorized(err)
	}
	ok, err := adminCtrl.srv.Update(&admin, &para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(ok)
}

func UpdateAdminPwd(ctx *kid.Context) any {
	if err := ctx.ValidFiled(kid.Rules{
		"uid": "required",
	}); err != nil {
		return web.ParamsErr(err)
	}
	uid := ctx.FindInt64("uid")
	var para dto.AdminUpdatePwdReq
	if err := ctx.Valid(&para); err != nil {
		return web.ParamsErr(err)
	}
	para.Uid = uid

	var admin models.Admin
	if err := ctx.GetBindUser(&admin); err != nil {
		return web.Unauthorized(err)
	}
	ok, err := adminCtrl.srv.UpdatePassword(&admin, &para)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(ok)
}

func DeleteAdmin(ctx *kid.Context) any {
	if err := ctx.ValidFiled(kid.Rules{
		"uid": "required",
	}); err != nil {
		return web.ParamsErr(err)
	}
	uid := ctx.FindInt64("uid")

	var admin models.Admin
	if err := ctx.GetBindUser(&admin); err != nil {
		return web.Unauthorized(err)
	}
	ok, err := adminCtrl.srv.Delete(&admin, uid)
	if err != nil {
		return web.Error(err)
	}
	return web.Success(ok)
}
