package middleware

import (
	"context"
	"github.com/leor-w/kid"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"
	"github.com/leor-w/kid/guard"
	"github.com/leor-w/kid/guard/jwt"
	"role_ai/infrastructure/ecode"
	"role_ai/infrastructure/web"
	"role_ai/models"
	"role_ai/repos"
)

type Authorize struct {
	adminRepo repos.IAdminRepository `inject:""`
}

var authorize *Authorize

func (auth *Authorize) Provide(_ context.Context) interface{} {
	authorize = &Authorize{}
	return authorize
}

type AuthUserKey struct{}

// AdminTokenAuth 后台管理 token 登录验证
func AdminTokenAuth(ctx *kid.Context) {
	user, err := GetUser(ctx)
	if err != nil {
		ctx.JSON(200, web.Error(err))
		ctx.Abort()
		return
	}
	if user.Type != guard.Admin {
		ctx.JSON(200, web.ErrorWithStatus(ecode.AuthFailedErr))
		ctx.Abort()
		return
	}
	var admin models.Admin
	err = authorize.adminRepo.GetOne(&finder.Finder{
		Wheres:    where.New().And(where.Eq("uid", user.Uid)),
		Recipient: &admin,
	})
	if err != nil {
		ctx.JSON(200, web.ErrorWithStatus(ecode.UserNotFoundErr))
		ctx.Abort()
		return
	}
	//// 获取管理员角色
	//var role models.Role
	//err = authorize.adminRepo.GetOne(&finder.Finder{
	//	Wheres:    where.New().And(where.Eq("id", admin.Role)),
	//	Recipient: &role,
	//})
	//if err != nil {
	//	ctx.JSON(200, web.ErrorWithStatus(ecode.RoleNotExistErr))
	//	ctx.Abort()
	//	return
	//}
	//// 如果对应角色为启用状态则设置其角色名.
	//// 如果为弃用状态不设置用户名. 后面验证权限时将会拒绝访问
	//if role.Enable == 1 {
	//	admin.RoleName = role.Name
	//}
	//// 记录token
	//admin.Token = ctx.FindString("token")
	//ctx.BindUser(&admin)
	ctx.Next()
}

//// AppTokenAuth 客户端登录 token 验证
//func AppTokenAuth(ctx *kid.Context) {
//	user, err := GetUser(ctx)
//	if err != nil {
//		ctx.JSON(200, web.Error(err))
//		ctx.Abort()
//		return
//	}
//	if user.Type != guard.General {
//		ctx.JSON(200, web.ErrorWithStatus(ecode.AuthFailedErr))
//		ctx.Abort()
//		return
//	}
//	var appUser models.User
//	err = authorize.userRepo.GetOneWithErr(
//		&repos.QueryWrapper{
//			Where:     mysql.Wrappers(mysql.NewWhereWrapper("uid", "=", user.Uid)),
//			Recipient: &appUser,
//		})
//	if err != nil {
//		ctx.JSON(200, web.ErrorWithStatus(ecode.UserNotFoundErr))
//		ctx.Abort()
//		return
//	}
//	ctx.BindUser(&appUser)
//	ctx.Next()
//}
//
//// AppTokenAuthIgnoreError 客户端登录 token 验证
//func AppTokenAuthIgnoreError(ctx *kid.Context) {
//	user, _ := GetUser(ctx)
//	if user != nil && user.Uid > 0 && user.Type == guard.General {
//		var appUser models.User
//		err := authorize.userRepo.GetOneWithErr(
//			&repos.QueryWrapper{
//				Where:     mysql.Wrappers(mysql.NewWhereWrapper("uid", "=", user.Uid)),
//				Recipient: &appUser,
//			})
//		if err != nil {
//			logger.Errorf("验证客户端用户错误: %s", err.Error())
//		}
//		ctx.BindUser(&appUser)
//	}
//	ctx.Next()
//}

func GetUser(ctx *kid.Context) (*guard.User, error) {
	tokenStr := ctx.FindString("token")
	if len(tokenStr) <= 0 {
		return nil, errors.New(ecode.NotFoundTokenErr)
	}
	user, err := jwt.Default().Verify(tokenStr)
	if err != nil {
		return nil, errors.New(ecode.TokenNotExistErr, err)
	}
	return user, nil
}
