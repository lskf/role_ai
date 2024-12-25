package repos

import (
	"context"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/repos"
)

type IRoleRepository interface {
	repos.IBasicRepository
}
type RoleRepository struct {
	*mysql.Repository `inject:""`
}

func (repo *RoleRepository) Provide(context.Context) any {
	return &RoleRepository{}
}
