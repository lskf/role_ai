package app

import (
	"fmt"
	"github.com/leor-w/kid/guard/jwt"
	"github.com/leor-w/kid/guard/store"
	"github.com/leor-w/kid/plugin/captcha"
	"github.com/leor-w/kid/plugin/qiniu"
	"role_ai/infrastructure/middleware"
	"role_ai/routers"

	"github.com/leor-w/injector"
	"github.com/leor-w/kid"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/redis"
	"github.com/leor-w/kid/plugin/lock"
	"role_ai/handler"
	"role_ai/models"
	"role_ai/repos"
	"role_ai/services"
)

func Run() {
	app := kid.New(
		kid.WithConfigs("./config/config.yml"),
		kid.WithRunMode("app.runMode"),
		kid.WithDbMigrate(&models.DatabaseMigrate{}),
		kid.WithAppName("RoleAi"),
	)
	ioc := injector.New()
	appScope := ioc.Scope("app")
	appScope.Provide(new(redis.Client))
	appScope.Provide(new(mysql.MySQL))
	appScope.Provide(new(jwt.Jwt))
	appScope.Provide(new(jwt.RedisBlacklist))
	appScope.Provide(new(store.RedisStore))
	appScope.Provide(new(lock.RedisLock))
	appScope.Provide(new(redis.RedisRepository))
	appScope.Provide(new(mysql.Repository))
	appScope.Provide(new(captcha.Captcha))
	appScope.Provide(new(captcha.RedisStore))
	appScope.Provide(new(middleware.Authorize))
	appScope.Provide(new(middleware.RequestSign))
	appScope.Provide(new(qiniu.Qiniu))

	handler.InitController(appScope)
	services.InitService(appScope)
	repos.InitRepository(appScope)

	if err := ioc.Populate(); err != nil {
		panic(fmt.Sprintf("依赖注入错误: %s", err.Error()))
	}
	routers.InitRouter(app)

	app.Launch(fmt.Sprintf("%s:%s", config.GetString("app.host"), config.GetString("app.port")))
}
