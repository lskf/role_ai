package repos

import (
	"context"
	"github.com/leor-w/injector"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/redis"
)

func InitRepository(appScope *injector.Scope) {
	appScope.Provide(new(UserRepository))
	appScope.Provide(new(AdminRepository))
	appScope.Provide(new(SignerRepository))
	appScope.Provide(new(SmsRepository))
	appScope.Provide(new(RoleRepository))
}

type BasisRepository struct {
	*mysql.Repository      `inject:""`
	*redis.RedisRepository `inject:""`
}

func (repo *BasisRepository) Provide(context.Context) interface{} {
	return repo
}
