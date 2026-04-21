package model

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user
type Ad struct {
	Id         int    `gorm:"primary_key;auto_increment;default();description(ID)"`
	Catid      int    `gorm:"default(0);null;description(分類id)"`
	Title      string `gorm:"size(50);default();null;index;description(标题)"`
	Img        string `gorm:"size(100);default();null;description(图片)"`
	Url        string `gorm:"size(100);default();null;description(链接)"`
	Sort       int16  `gorm:"default(0);null;description(排序: 数字小的排前面)"`
	Status     int8   `gorm:"default(1);null;description(状态: 0=禁用 1=启用)"`
	CreateTime int    `gorm:"default(0);null;description(創建時間)"`
	CreateUser string `gorm:"size(32);default(0);null;description(創建人)"`
	UpdateTime int    `gorm:"default(0);null;description(修改時間)"`
	UpdateUser string `gorm:"size(32);default(0);null;description(修改人)"`
	DeleteTime int    `gorm:"default(0);null;description(刪除時間)"`
	DeleteUser string `gorm:"size(32);default(0);null;description(刪除人)"`
}