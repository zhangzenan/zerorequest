// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"zerorequest/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	//初始化数据库连接
	dbClient, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		logx.Error(err)
		return nil
	}
	return &ServiceContext{
		Config: c,
		DB:     dbClient,
	}
}
