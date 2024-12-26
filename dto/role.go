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

type Voice struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
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
	Name string `json:"name"`
	Uid  int64  `json:"uid"`
	Sort int    `json:"sort"` // 1:按对话人数倒序，0：按创建时间倒序
	PageField
}

type RoleListResp struct {
	List  []Role `json:"list"`
	Total int64  `json:"total"`
}
