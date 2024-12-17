package repos

import (
	"context"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/redis"
	"github.com/leor-w/kid/database/repos"
)

type ILoginRepository interface {
	repos.IBasicRepository
	repos.IRedisRepository
}

type LoginRepository struct {
	*mysql.Repository      `inject:""`
	*redis.RedisRepository `inject:""`
}

func (repo *LoginRepository) Provide(context.Context) any {
	return &AdminRepository{}
}

func (repo *LoginRepository) GetCaptcha() {}
