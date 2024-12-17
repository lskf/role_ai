package constant

const (
	CodeTypeAdminLogin         = iota + 1 // 1 = 管理员登录
	CodeTypeClientLogin                   // 2 = 客户端用户登录
	CodeTypeOARegUserBindPhone            // 3 = 微信二维码登录绑定手机号
	CodeTypeUserVerified                  // 4 = 用户实名认证
	CodeTypeUpdatePhone                   // 5 = 修改手机号
)
