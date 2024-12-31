package common

const (
	// 运行环境
	RunModeRelease = "release"
	RunModeDebug   = "debug"
	RunModeTest    = "test"

	//时间格式
	TimeFormatToDateTime = "2006-01-02 15:04:05"
	TimeFormatToDate     = "2006-01-02"

	//上传文件
	UploadFileTypeUserAvatar  = 1 //用户头像
	UploadFileTypeRoleAvatar  = 2 //角色头像
	UploadFileTypeChatPicture = 3 //对话生成图片

	//上传文件路径
	UploadFilePathUserAvatar  = "user_avatar"  //用户头像
	UploadFilePathRoleAvatar  = "role_avatar"  //角色头像
	UploadFilePathChatPicture = "chat_picture" //对话生成图片
)
