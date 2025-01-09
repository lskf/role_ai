package repos

import (
	"context"
	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/repos"
	"role_ai/models"
)

type IRoleRepository interface {
	repos.IBasicRepository
	GetRoleByRoleNameAndChatUid(uid int64, roleName string) (list []models.Role, err error)
}
type RoleRepository struct {
	*mysql.Repository `inject:""`
}

func (repo *RoleRepository) Provide(context.Context) any {
	return &RoleRepository{}
}

func (repo *RoleRepository) GetRoleByRoleNameAndChatUid(uid int64, roleName string) (list []models.Role, err error) {
	db := repo.DB.Model(&models.Role{})
	db = db.Select("roles.id,roles.role_name,roles.avatar")
	db = db.Joins("LEFT JOIN chats ON roles.id = chats.role_id")
	db = db.Where("roles.role_name like ? AND chats.uid = ?", "%"+roleName+"%", uid)
	err = db.Find(&list).Error
	return list, err
}
