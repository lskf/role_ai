package models

const (
	ChatHistoryRoleUser      = 1
	ChatHistoryRoleAssistant = 2
	ChatHistoryTypeChat      = 1 //对话
	ReplyStyleLyrical        = 1
	ReplyStyleMaleNsfw       = 2
	ReplyStyleFemaleNsfw     = 3
)

// 聊天记录表
type Chat struct {
	IdField
	UidField
	RoleId       int64  `gorm:"column:role_id;type:int(11);default:0;comment:角色id;NOT NULL" json:"role_id"`
	ReplyStyle   int64  `gorm:"column:reply_style;type:tinyint(4);comment:回复风格，1：lyrical，2：male nsfw，3：female nsfw ;NOT NULL" json:"reply_style"`
	WordCount    int64  `gorm:"column:word_count;type:int(11);comment:字数规模;NOT NULL" json:"word_count"`
	Gamification string `gorm:"column:gamification;type:text;comment:游戏化" json:"gamification"`
	CreatedAtField
	UpdatedAtField
}

//type ChatHistory struct {
//	IdField
//	ChatId   int64  `gorm:"column:chat_id;type:int(11);default:0;comment:chat_id;NOT NULL" json:"chat_id"`
//	Question string `gorm:"column:question;type:text;comment:用户提问;NOT NULL" json:"question"`
//	Answer   string `gorm:"column:answer;type:text;comment:ai回复;NOT NULL" json:"answer"`
//	Abstract string `gorm:"column:abstract;type:varchar(500);comment:对话摘要;NOT NULL" json:"abstract"`
//	Content  string `gorm:"column:content;type:text;comment:完整内容;NOT NULL" json:"content"`
//	CreatedAtField
//	UpdatedAtField
//}

type ChatHistory struct {
	IdField
	ChatId   int64  `gorm:"column:chat_id;type:int(11);default:0;comment:chat_id;NOT NULL" json:"chat_id"`
	RoleType int64  `gorm:"column:role_type;type:tinyint(4);default:0;comment:角色类型 1：user，2：assistant;NOT NULL" json:"role_type"`
	Type     int64  `gorm:"column:type;type:tinyint(4);default:0;comment:1：聊天内容;NOT NULL" json:"type"`
	Content  string `gorm:"column:content;type:text;comment:内容" json:"content"`
	Info     string `gorm:"column:info;type:text" json:"info"`
	CreatedAtField
	UpdatedAtField
}

type Reply struct {
	Weekday   string `json:"weekday"`
	Time      string `json:"time"`
	Locations string `json:"locations"`
	Weather   string `json:"weather"`
	Content   string `json:"content"`
	Affection string `json:"affection"`
	Sexuality string `json:"sexuality"`
}
