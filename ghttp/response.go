package ghttp

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func Success(c *gin.Context, message string, data any) {
	if message == "" {
		message = "success"
	}
	res := gin.H{
		"code":    0,
		"data":    data,
		"message": message,
	}
	Json(c, &res)
}

func Fail(c *gin.Context, code int, message string, data any) {
	if message == "" {
		message = "fail"
	}
	res := gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	}
	Json(c, &res)
}

func Json(c *gin.Context, data *gin.H) {
	c.AsciiJSON(200, data)
}

func Write(c *gin.Context, data string) {
	c.Render(200, render.Data{Data: []byte(data)})
}

func WriteByte(c *gin.Context, data []byte) {
	c.Render(200, render.Data{Data: data})
}
