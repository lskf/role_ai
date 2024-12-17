package dto

import "regexp"

type CreatorAndUpdaterField struct {
	CreatorId   int64  `json:"creator_id"`
	UpdaterId   int64  `json:"updater_id"`
	CreatorName string `json:"creator_name"`
	UpdaterName string `json:"updater_name"`
}

func ParseCreatorUpdater(creatorId, updaterId int64, creatorName, updaterName string) CreatorAndUpdaterField {
	return CreatorAndUpdaterField{
		CreatorId:   creatorId,
		UpdaterId:   updaterId,
		CreatorName: creatorName,
		UpdaterName: updaterName,
	}
}

type CreatedUpdatedTimeField struct {
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func ParseCreatedUpdatedAt(created, updated int64) CreatedUpdatedTimeField {
	return CreatedUpdatedTimeField{
		CreatedAt: created,
		UpdatedAt: updated,
	}
}

type PageField struct {
	PageNum  int `json:"page_num" form:"page_num,default=1" `
	PageSize int `json:"page_size" form:"page_size,default=20"`
}

// IsDecimal 判断是否是数字
func IsDecimal(str string) bool {
	// 正则表达式模式
	pattern := `^\d+(\.\d+)?$`

	// 编译正则表达式
	regex := regexp.MustCompile(pattern)

	// 匹配字符串
	return regex.MatchString(str)
}
