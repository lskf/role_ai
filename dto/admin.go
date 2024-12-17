package dto

type Admin struct {
	Id       int64  `json:"-"`
	Uid      int64  `json:"uid"`       //用户id
	Phone    string `json:"phone"`     //手机号码
	Account  string `json:"account"`   //账号
	Avatar   string `json:"avatar"`    //头像
	Password string `json:"-"`         //密码
	Salt     string `json:"-"`         //密码盐
	Role     int64  `json:"role"`      //角色id
	RealName string `json:"real_name"` //真实姓名
	CreatorAndUpdaterField
	CreatedUpdatedTimeField
}

type AdminCreateReq struct {
	Phone    string `json:"phone" validate:"required"`    // 手机号
	Account  string `json:"account" validate:"required"`  // 账号
	Password string `json:"password" validate:"required"` // 密码
	Avatar   string `json:"avatar"`                       // 头像
	RealName string `json:"real_name"`                    // 真实姓名
	Role     int64  `json:"role"`                         // 角色id
}

type AdminUpdateReq struct {
	Uid      int64  `json:"uid"`
	Phone    string `json:"phone" validate:"required"`   // 手机号
	Account  string `json:"account" validate:"required"` // 账号
	Avatar   string `json:"avatar"`                      // 头像
	RealName string `json:"real_name"`                   // 真实姓名
	Role     int64  `json:"role"`                        // 角色id
}

type AdminUpdatePwdReq struct {
	Uid    int64  `json:"uid"`                         //uid
	NewPwd string `json:"new_pwd" validate:"required"` // 密码
}

type AdminListReq struct {
	Uid      int64  `json:"uid" form:"uid"`
	Phone    string `json:"phone" form:"phone"`
	Account  string `json:"account" form:"account"`
	Role     int64  `json:"role" form:"role"`
	RealName string `json:"real_name" form:"real_name"`
	PageField
}

type AdminListResp struct {
	List  []Admin `json:"list"`
	Total int64   `json:"total"`
}
