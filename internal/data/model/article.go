package model

//数据库结构, 这里只写gorm tag 不用json tag
type Article struct {
	Id         int    `gorm:"primaryKey;autoIncrement;comment:ID"`
	Catid      int    `gorm:"default:0;comment:分类id"`
	Title      string `gorm:"size:50;index;comment:标题"`
	Info       string `gorm:"comment:简介"`
	Content    string `gorm:"type:text;comment:内容"`
	Img        string `gorm:"size:100;comment:图片"`
	Author     string `gorm:"size:30;index;comment:作者"`
	Extra      string `gorm:"comment:扩展"`
	Sort       int16  `gorm:"default:0;comment:排序"`
	Status     int8   `gorm:"default:1;comment:状态"`
	CreateTime int    `gorm:"comment:创建时间"`
	CreateUser string `gorm:"size:32;comment:创建人"`
	UpdateTime int    `gorm:"comment:更新时间"`
	UpdateUser string `gorm:"size:32;comment:更新人"`
	DeleteTime int    `gorm:"index;comment:删除时间"`
	DeleteUser string `gorm:"size:32;comment:删除人"`
}