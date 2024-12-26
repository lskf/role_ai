package models

type Role struct {
	IdField
	UidField
	Avatar       string `gorm:"column:avatar;type:varchar(255);comment:头像;NOT NULL" json:"avatar"`
	RoleName     string `gorm:"column:role_name;type:varchar(255);comment:角色名称;NOT NULL" json:"role_name"`
	Gender       string `gorm:"column:gender;type:varchar(50);comment:性别;NOT NULL" json:"gender"`
	Desc         string `gorm:"column:desc;type:text;comment:角色简介;NOT NULL" json:"desc"`
	Worldview    string `gorm:"column:worldview;type:text;comment:世界观;NOT NULL" json:"worldview"`
	Remark       string `gorm:"column:remark;type:varchar(255);comment:开场白;NOT NULL" json:"remark"`
	Tag          string `gorm:"column:tag;type:varchar(500);comment:标签;NOT NULL" json:"tag"`
	Gamification string `gorm:"column:gamification;type:text;comment:游戏化;NOT NULL" json:"gamification"`
	IsPublic     int64  `gorm:"column:is_public;type:tinyint(1);default:0;comment:是否公开 1：公开，2：私密;NOT NULL" json:"is_public"`
	VoiceId      int64  `gorm:"column:voice_id;type:int(11);default:0;comment:声音模型id;NOT NULL" json:"voice_id"`
	ChatNum      int64  `gorm:"column:chat_num;type:int(11);default:0;comment:对话人数;NOT NULL" json:"chat_num"`
	CreatedAtField
	UpdatedAtField
}

type RoleStyle struct {
	IdField
	RoleId  int64  `gorm:"column:role_id;type:int(11);default:0;comment:角色id;NOT NULL" json:"role_id"`
	Content string `gorm:"column:content;type:text;comment:内容;NOT NULL" json:"content"`
}

const (
	PubilcRole  = 1
	PrivateRole = 2
)
