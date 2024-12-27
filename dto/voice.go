package dto

type Voice struct {
	Id      int64    `json:"id"`
	Name    string   `json:"name"`
	Desc    string   `json:"-"`
	Sample  string   `json:"sample"`
	Pth     string   `json:"-"`
	Ckpt    string   `json:"-"`
	DescArr []string `json:"desc"`
}

type VoiceListReq struct {
	PageField
}

type VoiceListResp struct {
	List  []Voice `json:"list"`
	Total int64   `json:"total"`
}
