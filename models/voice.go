package models

type Voice struct {
	IdField
	Name   string `gorm:"column:name;type:varchar(255);comment:声音名称;NOT NULL" json:"name"`
	Desc   string `gorm:"column:desc;type:varchar(255);comment:简介;NOT NULL" json:"desc"`
	Sample string `gorm:"column:sample;type:varchar(500);comment:声音样本文件路径;NOT NULL" json:"sample"`
	Pth    string `gorm:"column:pth;type:varchar(255);comment:pth文件路径;NOT NULL" json:"pth"`
	Ckpt   string `gorm:"column:ckpt;type:varchar(255);comment:ckpt文件路径;NOT NULL" json:"ckpt"`
	CreatedAtField
	UpdatedAtField
}
