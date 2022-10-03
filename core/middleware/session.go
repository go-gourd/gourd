package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	gdsession "github.com/go-gourd/gourd/core/sessions"
	"github.com/go-gourd/gourd/log"
	"github.com/go-gourd/gourd/utils"
)

// SessionMiddle 获取SESSION中间件
func SessionMiddle() gin.HandlerFunc {

	cfg := gdsession.SessionConfig{}

	err := config.ParseConfig("session", &cfg)
	if err != nil {
		log.Error(err.Error())
	}

	//创建目录
	err = utils.CheckAndMkdir(cfg.Path)
	if err != nil {
		log.Error(err.Error())
	}

	if cfg.Type == "file" {
		store := gdsession.NewFileStore(cfg.Path, []byte("gourd"))

		store.Options(sessions.Options{
			MaxAge: cfg.Expire, //过期时间，秒
			Domain: cfg.Domain,
			Secure: cfg.Secure,
		})

		// 设置Session中间件，参数1是SessionName，参数2是储存引擎
		return sessions.Sessions("session", store)
	} else if cfg.Type == "cookie" {
		// 创建基于 cookie 的存储引擎
		store := cookie.NewStore([]byte("gourd"))

		// 设置Session中间件，参数1是SessionName，参数2是储存引擎
		return sessions.Sessions("session", store)
	} else {
		return nil
	}
}
