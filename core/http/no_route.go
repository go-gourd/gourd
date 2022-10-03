package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	"os"
)

func InitDefaultRoute(router *gin.Engine) {

	//注册默认路由
	router.NoRoute(func(c *gin.Context) {

		//App配置获取
		var cfg config.AppConfig
		err := config.ParseConfig("app", &cfg)

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
			}
		}

		//404操作

	})
}
