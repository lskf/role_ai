package dto

type Chat struct {
	Id         int64         `json:"id"`
	RoleDetail Role          `json:"role_detail"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
	Histories  []ChatHistory `json:"histories"`
}

type ChatHistory struct {
	Id        int64  `json:"id"`
	ChatId    int64  `json:"chat_id"`
	RoleType  int64  `json:"role_type"` //1:User, 2:Assistant
	Type      int64  `json:"type"`      //对话内容
	Content   string `json:"content"`   //内容
	Affection int64  `json:"affection"`
	Sexuality int64  `json:"sexuality"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

type ChatReq struct {
	RoleId   int64  `json:"role_id"`
	Question string `json:"question"`
}

type ChatResp struct{}

type ChatListReq struct {
	Name string `json:"name" form:"name"`
	PageField
}

type ChatListResp struct {
	List  []Chat `json:"list"`
	Total int64  `json:"total"`
}
