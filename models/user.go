package models

import "time"

type User struct {
	IdField
	UidField
	Account  string    `gorm:"column:account;type:varchar(20);comment:账号;NOT NULL" json:"account"`
	Password string    `gorm:"column:password;type:varchar(50);comment:密码;NOT NULL" json:"password"`
	Phone    string    `gorm:"column:phone;type:varchar(20);comment:手机号码;NOT NULL" json:"phone"`
	NickName string    `gorm:"column:nick_name;type:varchar(30);comment:昵称;NOT NULL" json:"nick_name"`
	Avatar   string    `gorm:"column:avatar;type:varchar(30);comment:头像;NOT NULL" json:"avatar"`
	Gender   int64     `gorm:"column:gender;type:tinyint(4);default:0;comment:性别 1：男，2：女;NOT NULL" json:"gender"`
	Birthday time.Time `gorm:"column:birthday;type:datetime;default:CURRENT_TIMESTAMP;comment:生日日期;NOT NULL" json:"birthday"`
	Menber   int64     `gorm:"column:menber;type:tinyint(4);default:0;comment:会员 1：普通用户，2：初级会员，3：中级会员，4：高级会员;NOT NULL" json:"menber"`
	Status   int64     `gorm:"column:status;type:tinyint(4);default:0;comment:状态 1：正常，2：禁用;NOT NULL" json:"status"`
	CreatedAtField
	UpdatedAtField
	Token string `gorm:"-"`
}

const (
	//会员等级
	UserMenberNormal = iota + 1 //普通用户
	UserMenberNovice            //初级会员
	UserMenberMiddle            //中级会员
	UserMenberHigh              //高级会员
)
