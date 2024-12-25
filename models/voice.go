package models

type Voice struct {
	IdField
	Name string `gorm:"column:name;type:varchar(255);comment:声音名称;NOT NULL" json:"name"`
	Desc string `gorm:"column:desc;type:varchar(255);comment:简介;NOT NULL" json:"desc"`
	CreatedAtField
	UpdatedAtField
}
