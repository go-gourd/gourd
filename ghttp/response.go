package ghttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/go-gourd/gourd/config"
	"os"
)

func Write(c *gin.Context, data string) {
	c.Render(200, render.Data{Data: []byte(data)})
}

func WriteByte(c *gin.Context, data []byte) {
	c.Render(200, render.Data{Data: data})
}

func WriteStaticFile(c *gin.Context, filepath string) error {
	conf := config.GetHttpConfig()
	if filepath == "" {
		filepath = conf.Public + c.Request.URL.Path
	}
	//判断文件是否存在
	_, err := os.Stat(filepath)
	if err == nil {
		return errors.New("file not exist: " + filepath)
	}

	c.File(filepath)
	return nil
}
