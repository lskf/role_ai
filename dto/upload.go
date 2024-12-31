package dto

type UploadFileResp struct {
	FilePath string `json:"file_path"`
}

type UploadFileByUrlReq struct {
	Url          string `json:"url"`
	FilePathType int64  `json:"file_path_type"`
}
