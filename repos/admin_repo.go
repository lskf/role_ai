package repos

import (
	"context"
	"github.com/leor-w/kid/database/redis"
	"role_ai/infrastructure/tools"
	"time"

	"github.com/leor-w/kid/database/mysql"
	"github.com/leor-w/kid/database/repos"
	"github.com/leor-w/kid/database/repos/finder"
	"github.com/leor-w/kid/database/repos/where"
	"github.com/leor-w/kid/errors"

	"role_ai/infrastructure/constant"
	"role_ai/infrastructure/ecode"
	"role_ai/models"
)

type IAdminRepository interface {
	repos.IBasicRepository
	repos.IRedisRepository
	CreateAdmin(admin *models.Admin) error

	FindAllAdmins(query *finder.Finder) error
	PhoneOrAccountExisted(phone, account string) (existed bool, err error)
	SaveLoginAuthCode(phone, authCode string) error
	GetLoginAuthCode(phone string) (string, error)
	DiscardAuthCode(phone string) error
	GenerateUid() int64
}

type AdminRepository struct {
	*mysql.Repository      `inject:""`
	*redis.RedisRepository `inject:""`
}

func (repo *AdminRepository) Provide(context.Context) any {
	return &AdminRepository{}
}

func (repo *AdminRepository) CreateAdmin(admin *models.Admin) error {
	adminId := repo.GetUniqueID(&finder.Finder{
		Model:  new(models.Admin),
		Wheres: where.New().And(where.Eq("uid", nil)),
	}, 10000000, 99999999, 0, 0)
	admin.AdminId = adminId
	if err := repo.DB.Create(&admin).Error; err != nil {
		return errors.New(ecode.DatabaseErr, err)
	}
	return nil
}

func (repo *AdminRepository) FindAllAdmins(query *finder.Finder) error {
	return repo.DB.Model(&models.Admin{}).
		Joins("LEFT JOIN roles ON admins.role = roles.id").
		Select("admins.id", "admins.uid", "admins.phone",
			"admins.account", "admins.avatar", "admins.password",
			"admins.creator_id", "roles.name", "admins.created_at",
			"admins.updated_at", "admins.real_name", "admins.role").
		Scopes(mysql.Paginate(query.Num, query.Size)).
		Find(query.Recipient).Error
}

func (repo *AdminRepository) PhoneOrAccountExisted(phone, account string) (existed bool, err error) {
	var count int64
	existed = true //默认为存在
	db := repo.DB.Model(&models.Admin{})
	db = db.Where("(phone = ? or account = ?)", phone, account)
	err = db.Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (repo *AdminRepository) SaveLoginAuthCode(phone, authCode string) error {
	key := constant.GetAdminLoginAuthCodeKey(phone)
	if err := repo.RDB.Set(key, authCode, time.Minute*5).Err(); err != nil {
		return errors.New(ecode.DatabaseErr, err)
	}
	return nil
}

func (repo *AdminRepository) GetLoginAuthCode(phone string) (string, error) {
	key := constant.GetAdminLoginAuthCodeKey(phone)
	code, err := repo.RDB.Get(key).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

func (repo *AdminRepository) DiscardAuthCode(phone string) error {
	key := constant.GetAdminLoginAuthCodeKey(phone)
	return repo.RDB.Expire(key, 0).Err()
}

// GenerateUid 获取唯一的uid
func (repo *AdminRepository) GenerateUid() int64 {
	uid := tools.RandomInt64InRange(100000, 999999)
	for {
		exist := repo.Exist(&finder.Finder{
			Model:  &models.Admin{},
			Wheres: where.New().And(where.Eq("uid", uid)),
		})
		if exist {
			continue
		}
		break
	}
	return uid
}
