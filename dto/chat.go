package dto

type Chat struct {
	Id         int64          `json:"id"`
	RoleId     int64          `json:"role_id"`
	RoleName   string         `json:"role_name"`
	RoleAvatar string         `json:"role_avatar"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
	Histories  []*ChatHistory `json:"histories"`
}

type ChatHistory struct {
	Id        int64  `json:"id"`
	RoleType  int64  `json:"role_type"` //1:User, 2:Assistant
	Type      int64  `json:"type"`      //对话内容
	Content   string `json:"content"`   //内容
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

type ChatReq struct {
	RoleId   int64  `json:"role_id"`
	Question string `json:"question"`
}

type ChatResp struct {
	ChatId          int64         `json:"chat_id"`
	Affection       int64         `json:"affection"`
	Sexuality       int64         `json:"sexuality"`
	CreatedAt       string        `json:"-"`
	UpdatedAt       string        `json:"-"`
	ChatHistoryList []ChatHistory `json:"chat_history"`
}

type ChatListReq struct {
	Name string `json:"name" form:"name"`
	PageField
}

type ChatListResp struct {
	List  []Chat `json:"list"`
	Total int64  `json:"total"`
}

type ChatHistoryListReq struct {
	RoleId int64 `json:"role_id" form:"role_id"`
	ChatId int64 `json:"chat_id" form:"chat_id"`
	Id     int64 `json:"id" form:"id"`
	PageField
}

type ChatHistoryListResp struct {
	List  []ChatHistory `json:"list"`
	Total int64         `json:"total"`
}

type ChatTtsReq struct {
	Content string `json:"content" validate:"required"`
}
