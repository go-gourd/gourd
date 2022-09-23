package gourd

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

func (con BaseController) Success(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
	})
}

func (con BaseController) Error(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 1,
		"msg":  msg,
	})
}
