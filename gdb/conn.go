package gdb

import (
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// GetMysqlDb 获取Mysql连接
func GetMysqlDb() *gorm.DB {

	if db != nil {
		return db
	}

	conf := config.GetDbConfig()

	dsnParam := ""
	if conf.Param != "" {
		dsnParam = "?" + conf.Param
	}
	dsnF := "%s:%s@(%s:%d)/%s%s"
	dsn := fmt.Sprintf(dsnF, conf.User, conf.Pass, conf.Host, conf.Port, conf.Database, dsnParam)

	// 连接数据库
	newDb, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Error("cannot establish db connection: %w" + err.Error())
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	db = newDb

	return db
}
