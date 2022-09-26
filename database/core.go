package database

import (
	"github.com/go-gourd/gourd/config"
	glog "github.com/go-gourd/gourd/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

var db *gorm.DB

func ConnDb() {

	//App配置获取
	var cfg config.Database
	err := config.GetConfig("database", &cfg)
	if err != nil {
		glog.Info(err.Error())
	}

	dsn := cfg.User + ":" + cfg.Pass + "@tcp(" + cfg.Host + ":" + strconv.Itoa(int(cfg.Port)) +
		")/" + cfg.Db + "?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	//根据*grom.DB对象获得*sql.DB的通用数据库接口
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	sqlDB, _ := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

}

// GetDb 获取数据库连接
func GetDb() *gorm.DB {
	if db == nil {
		ConnDb()
	}
	return db
}
