package dto

type VoiceListReq struct {
	PageField
}

type VoiceListResp struct {
	List  []Voice `json:"list"`
	Total int64   `json:"total"`
}
