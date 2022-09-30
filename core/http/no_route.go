package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/log"
	"os"
)

func InitDefaultRoute(router *gin.Engine) {

	//注册默认路由
	router.NoRoute(func(c *gin.Context) {

		//App配置获取
		var cfg config.AppConfig
		err := config.GetConfig("app", &cfg)

		log.Info(cfg.ReleaseMode)

		if err == nil {
			if cfg.PublicPath != "" {
				//输出静态文件
				fileName := cfg.PublicPath + c.Request.URL.String()

				//判断静态文件是否存在
				_, err = os.Stat(fileName)
				if err == nil {
					//响应文件
					c.File(fileName)
					return
				}
			} else {
				log.Info("路径配置：" + cfg.PublicPath)
			}
		} else {
			log.Error(err.Error())
		}

		//404操作

	})
}
