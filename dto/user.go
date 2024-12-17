package dto

import "time"

type User struct {
	Id        int64     `json:"-"`
	Uid       int64     `json:"uid"`      //用户id
	Account   string    `json:"account"`  //账号
	Password  string    `json:"-"`        //密码
	Phone     string    `json:"phone"`    //手机号码
	NickName  string    `json:"nickname"` //昵称
	Avatar    string    `json:"avatar"`   //头像
	Gender    int64     `json:"gender"`   //性别 1：男，2：女
	Birthday  time.Time `json:"birthday"` //生日
	Menber    int64     `json:"menber"`   //会员 1：普通用户，2：初级会员，3：中级会员，4：高级会员
	Status    int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
