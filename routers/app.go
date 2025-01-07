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

	logged := group.Group("", middleware.AppTokenAuth) // 需要登录的路由
	{
		logout := logged.Group("/logout") // logout 的路由
		{
			logout.POST("", handler.Logout) // 退出登录
		}

		user := logged.Group("/user")
		{
			user.GET("", handler.GetUserDetail) //获取用户详情
			user.PUT("", handler.UpdateUser)    //更新用户信息
		}

		role := logged.Group("/role")
		{
			role.GET("/list", handler.GetRoleList)
			role.GET("/:id", handler.GetRoleDetail)
			role.POST("", handler.CreateRole)
			role.PUT("/:id", handler.UpdateRole)
			role.POST("/getRoleAvatarSetting", handler.GetRoleAvatarSetting)
			role.POST("/getRoleSetting", handler.GetRoleSetting)
			role.POST("/createRoleAvatar", handler.CreateRoleAvatar)
			role.GET("/getRoleAvatarHistory/:prompt_id", handler.GetRoleAvatarHistory)
			role.GET("/getRoleAvatar", handler.GetRoleAvatar)
		}

		voice := logged.Group("/voice")
		{
			voice.GET("", handler.GetVoiceList)
		}

		upload := logged.Group("")
		{
			upload.POST("/upload", handler.Upload)
			upload.POST("/uploadByUrl", handler.UploadByUrl)
		}

		chat := logged.Group("/chat")
		{
			chat.POST("", handler.Chat)
		}
	}

}
