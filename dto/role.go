package dto

import "time"

type Role struct {
	Id              int64            `json:"id"`
	Uid             int64            `json:"uid"`                           //用户id
	Avatar          string           `json:"avatar" validate:"required"`    //头像
	RoleName        string           `json:"role_name" validate:"required"` //角色名称
	Gender          string           `json:"gender" validate:"required"`    //性别
	Desc            string           `json:"desc" validate:"required"`      //角色简介
	Worldview       string           `json:"worldview" validate:"required"` //世界观
	Remark          string           `json:"remark" validate:"required"`    //开场白
	Tag             string           `json:"-"`                             //标签
	Gamification    string           `json:"-"`                             //游戏化
	IsPublic        int64            `json:"is_public" validate:"required"` //是否公开 1：公开，2：私密
	VoiceId         int64            `json:"voice_id"`                      //声音模型id
	ChatNum         int64            `json:"chat_num"`                      //对话人数
	CreatedAt       time.Time        `json:"-"`
	UpdatedAt       time.Time        `json:"-"`
	TagArray        []string         `json:"tag"`
	StyleArray      []SpeechStyleObj `json:"style"`
	GamificationObj `json:"gamification"`
	CreatedAtStr    string `json:"created_at"`
	UpdatedAtStr    string `json:"updated_at"`
}

type RoleStyle struct {
	Id      int64  `json:"id"`
	RoleId  int    `json:"role_id"`
	Content string `json:"content"`
}

type SpeechStyleObj struct {
	User      string `json:"user"`      //用户输入
	Assistant string `json:"assistant"` //角色回答
}

type GamificationObj struct {
	AffectionInitial int64 `json:"affection_initial"`
	SexualityInitial int64 `json:"sexuality_initial"`
}

type CreateRoleReq struct {
	Role
}

type UpdateRoleResp struct {
	Id              int64            `json:"id"`                            //Id
	Avatar          string           `json:"avatar" validate:"required"`    //头像
	RoleName        string           `json:"role_name" validate:"required"` //角色名称
	Gender          string           `json:"gender" validate:"required"`    //性别
	Desc            string           `json:"desc" validate:"required"`      //角色简介
	Worldview       string           `json:"worldview" validate:"required"` //世界观
	Remark          string           `json:"remark" validate:"required"`    //开场白
	Tag             string           `json:"-"`                             //标签
	Gamification    string           `json:"-"`                             //游戏化
	IsPublic        int64            `json:"is_public" validate:"required"` //是否公开 1：公开，2：私密
	VoiceId         int64            `json:"voice_id"`                      //声音模型id
	TagArray        []string         `json:"tag"`
	StyleArray      []SpeechStyleObj `json:"style"`
	GamificationObj `json:"gamification"`
}

type RoleListReq struct {
	Name string `json:"name" form:"name"`
	Uid  int64  `json:"uid" form:"uid"`
	Sort int    `json:"sort" json:"sort"` // 1:按对话人数倒序，0：按创建时间倒序
	PageField
}

type RoleListResp struct {
	List  []Role `json:"list"`
	Total int64  `json:"total"`
}

type AiCreateRoleReq struct {
	Gender        []string `json:"gender"`
	StoryGenre    []string `json:"story_genre"`
	RoleType      []string `json:"role_type"`
	Personality   []string `json:"personality"`
	Interests     []string `json:"interests"`
	Preferences   []string `json:"preferences"`
	Dislike       []string `json:"dislike"`
	Background    []string `json:"background"`
	Relationships []string `json:"relationships"`
	Quirks        []string `json:"quirks"`
	ArtStyle      string   `json:"art_style"`
}

type GetRoleAvatarResq struct {
	ArtStyle string `json:"art_style"` //画风
	Desc     string `json:"desc"`      //描述
}
