package ghttp

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
)

type routerHandler func(*gin.Engine)

var ginEngine *gin.Engine

// GetEngine 获取Gin容器
func GetEngine() *gin.Engine {

	if ginEngine != nil {
		return ginEngine
	}

	if config.GetAppConfig().Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginEngine = gin.New()

	return ginEngine
}

// SetRouter 设置路由
func SetRouter(call routerHandler) {
	call(GetEngine())
}
