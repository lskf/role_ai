package models

import "role_ai/infrastructure/utils"

type Admin struct {
	IdField
	AdminIdField
	Phone    string `gorm:"column:phone;type:varchar(50);uniqueIndex:idx_phone;comment:手机号"`
	Account  string `gorm:"column:account;type:varchar(50);uniqueIndex:idx_account;comment:账号"`
	Avatar   string `gorm:"column:avatar;type:varchar(256);comment:头像"`
	Password string `gorm:"column:password;type:varchar(64);comment:密码"`
	Salt     string `gorm:"column:salt;type:varchar(64);comment:密码盐"`
	Role     int64  `gorm:"column:role;type:bigint;comment:角色"`
	RealName string `gorm:"column:real_name;type:varchar(32);default:'';comment:真实姓名"`
	CreatorField
	UpdaterField
	DeletedFiled
	RoleName string `gorm:"-"`
	Token    string `gorm:"-"`
}

// CheckPassword 检查密码是否正确
func (a *Admin) CheckPassword(pwd string) bool {
	return utils.VerifyPassword(pwd, a.Salt, a.Password)
}

// EncodePassword 加密密码
func (a *Admin) EncodePassword() error {
	var err error
	a.Salt, err = utils.GenerateSalt()
	if err != nil {
		return err
	}
	a.Password = utils.Encode(a.Password, a.Salt)
	return nil
}
