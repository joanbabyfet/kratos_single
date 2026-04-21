package model

// 定义结构体, 字段首字母要大写才能进行json解析, 会自动转蛇底命令例 create_user, unique唯一索引
type Admin struct {
	Id           string `gorm:"primary_key;size(32);default();description(ID)"`
	Username     string `gorm:"unique;size(40);default();null;index;description(帐号)"`
	Password     string `gorm:"size(60);default();null;description(密码)"` //密码不输出, 改为可以接收
	Realname     string `gorm:"size(50);default();null;index;description(姓名)"`
	Sex          int8   `gorm:"default(1);null;description(性别 0=女 1=男)"`
	Email        string `gorm:"unique;size(100);default();null;index;description(信箱)"`
	Salt         string `gorm:"size(128);default();null;description(加密钥匙)"`
	RoleId       int    `gorm:"default(0);null;description(角色)"`
	RegIp        string `gorm:"size(15);default();null;description(注册ip)"`
	LoginTime    int    `gorm:"default(0);null;description(最后登录时间)"`
	LoginIp      string `gorm:"size(15);default();null;description(最后登录IP)"`
	LoginCountry string `gorm:"size(2);default();null;description(最后登录国家)"`
	Status       int8   `gorm:"default(1);null;description(状态: 0=禁用 1=启用)"`
	CreateTime   int    `gorm:"default(0);null;description(創建時間)"`
	CreateUser   string `gorm:"size(32);default(0);null;description(創建人)"`
	UpdateTime   int    `gorm:"default(0);null;description(修改時間)"`
	UpdateUser   string `gorm:"size(32);default(0);null;description(修改人)"`
	DeleteTime   int    `gorm:"default(0);null;description(刪除時間)"`
	DeleteUser   string `gorm:"size(32);default(0);null;description(刪除人)"`
}