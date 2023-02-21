package gdb

import (
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

var db *gorm.DB

type LogWriter struct{}

func (w LogWriter) Printf(format string, args ...any) {
	logger.Errorf(format, args...)
}

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

	// 慢日志阈值
	slowLogTime := conf.SlowLogTime
	if slowLogTime == 0 {
		slowLogTime = 60000 //默认1分钟
	}

	newLogger := gormLogger.New(
		LogWriter{},
		gormLogger.Config{
			SlowThreshold:             time.Duration(slowLogTime) * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  gormLogger.Warn,                               // 日志级别
			IgnoreRecordNotFoundError: true,                                          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                                         // 禁用彩色打印
		},
	)

	// 连接数据库
	newDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logger.Error("cannot establish db connection: %w" + err.Error())
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	db = newDb

	return db
}
