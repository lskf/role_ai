package handler

import (
	"github.com/leor-w/injector"
)

func InitController(appScope *injector.Scope) {
	appScope.Provide(new(LoginController))
	appScope.Provide(new(AdminController))
	appScope.Provide(new(UploadController))
	appScope.Provide(new(UserController))
	appScope.Provide(new(RoleController))
}
