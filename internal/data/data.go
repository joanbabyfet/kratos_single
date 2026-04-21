package data

import (
	"context"
	"kratos_single/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedisClient, NewArticleRepo, NewUserRepo, NewAdRepo, NewAdminRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB //加数据库
}

//使用事务
func (d *Data) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	// 初始化 GORM + 表前缀, gorm.Config 里只能用 gorm 的 logger（gormLogger）
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.Database.Prefix, // 前缀从配置读取
			SingularTable: true,              // 表名不加 s（推荐）
		},
	})
	if err != nil {
		return nil, nil, err
	}
	
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}
