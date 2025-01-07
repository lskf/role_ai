package dto

type Chat struct {
	Id         int64         `json:"id"`
	RoleDetail Role          `json:"role_detail"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
	Histories  []ChatHistory `json:"histories"`
}

type ChatHistory struct {
	Id       int64  `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
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
