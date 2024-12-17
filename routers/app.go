package routers

import (
	"github.com/leor-w/kid"
	"role_ai/handler"
	"role_ai/infrastructure/middleware"
)

func InitAppRouter(group *kid.RouterGroup) {
	login := group.Group("/login") // login 的路由
	{
		login.GET("/captcha", handler.GetLoginCaptcha)   // 获取登录验证码
		login.POST("/sendSmsCode", handler.SendSmsCode)  // 发送短信验证码
		login.POST("/loginBySms", handler.LoginBySms)    // 短信验证码登录
		login.POST("/password", handler.LoginByPassword) // 密码登录
	}

	logged := group.Group("", middleware.AdminTokenAuth) // 需要登录的路由
	{
		logout := logged.Group("/logout") // logout 的路由
		{
			logout.POST("", handler.Logout) // 退出登录
		}

		admin := logged.Group("/admin") // app 的路由
		{
			admin.POST("", handler.CreateAdmin)                  // 创建 admin
			admin.GET("", handler.GetAdminList)                  // 获取 admin 列表
			admin.GET("/:uid", handler.GetAdminDetail)           // 获取 admin 详情
			admin.PUT("/:uid", handler.UpdateAdmin)              // 更新 admin
			admin.PUT("/updatePwd/:uid", handler.UpdateAdminPwd) // 更新 admin 密码
			admin.DELETE("/:uid", handler.DeleteAdmin)           // 删除 admin
		}

		upload := logged.Group("/upload")
		{
			upload.GET("/token", handler.GetUploadToken)
		}
	}

}
