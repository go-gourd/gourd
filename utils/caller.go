package utils

import (
	"fmt"
	"path"
	"runtime"
)

func GetCallerInfo(skip int) (fileName, funcName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip + 1)
	if !ok {
		fmt.Println("runtime.Caller() failed")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file) // Base函数返回路径的最后一个元素
	return
}
