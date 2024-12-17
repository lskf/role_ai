package repos

import (
	"context"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/repos"
)

type IUserRepository interface {
	repos.IBasicRepository
}
type UserRepository struct {
	*mysql.Repository `inject:""`
}

func (repo *UserRepository) Provide(context.Context) any {
	return &UserRepository{}
}
