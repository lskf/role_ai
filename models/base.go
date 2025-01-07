package models

import (
	"gorm.io/plugin/soft_delete"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// IdField id主键
type IdField struct {
	Id int64 `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:id" json:"-"`
}

// UidField 用户id
type UidField struct {
	Uid int64 `gorm:"column:uid;type:int(11);default:0;comment:用户id;NOT NULL" json:"uid"`
}

// AdminIdField 管理员id
type AdminIdField struct {
	AdminId int64 `gorm:"column:admin_id;type:int(11);default:0;comment:管理员id;NOT NULL" json:"admin_id"`
}

// CreatedAtField 创建时间
type CreatedAtField struct {
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"created_at"`
}

func (model *CreatedAtField) BeforeCreate(*gorm.DB) error {
	model.CreatedAt = time.Now()
	return nil
}

// UpdatedAtField 更新时间
type UpdatedAtField struct {
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}

func (model *UpdatedAtField) BeforeCreate(*gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}

func (model *UpdatedAtField) BeforeUpdate(*gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}

// DeletedFiled 软删除
type DeletedFiled struct {
	Deleted soft_delete.DeletedAt `gorm:"softDelete:flag;default:0;"`
}

type CreatorField struct {
	CreatorId   int64  `gorm:"column:creator_id;type:bigint;default:0;comment:创建者"`
	CreatorName string `gorm:"column:creator_name;type:varchar(32);default:'';comment:创建者名字"`
}

func (c *CreatorField) SetCreator(uid int64, creatorName string) {
	c.CreatorId = uid
	c.CreatorName = creatorName
}

type UpdaterField struct {
	UpdaterId   int64  `gorm:"column:updater_id;type:bigint;default:0;comment:修改者"`
	UpdaterName string `gorm:"column:updater_name;type:varchar(32);default:'';comment:修改者名字"`
}

func (u *UpdaterField) SetUpdater(uid int64, updateName string) {
	u.UpdaterId = uid
	u.UpdaterName = updateName
}

type DeleterField struct {
	DeleterId  int64  `gorm:"column:deleter_id;type:bigint;default:0;comment:删除者"`
	DeleteName string `gorm:"column:deleter_name;type:varchar(32);default:'';comment:删除者名字"`
}

func (d *DeleterField) SetDeleter(uid int64, deleteName string) {
	d.DeleterId = uid
	d.DeleteName = deleteName
}

// DatabaseMigrate 数据库迁移接口实现
type DatabaseMigrate struct{}

// Models 在此处注册返回需要自动迁移的模型
func (dm *DatabaseMigrate) Models() []interface{} {
	return []interface{}{
		new(Admin), new(User), new(Role), new(RoleStyle), new(Voice), new(Chat), new(ChatHistory),
	}
}

type ExtraField struct {
	Extra string `gorm:"column:extra;type:text;comment:扩展"`
}

func Copy(dest, src interface{}) error {
	return copier.Copy(dest, src)
}

func GetDeletedMark() string {
	return "del" + strconv.FormatInt(time.Now().Unix(), 10)
}

// VerifyCode 验证码
type VerifyCode struct {
	Code    string `json:"code"`    //验证码
	Issued  int64  `json:"issued"`  //生效时间 时间戳
	Expired int64  `json:"expired"` //过期时间 时间戳
}
