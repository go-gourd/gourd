package test

import (
	"fmt"
	"github.com/go-gourd/gourd/cmd"
	"testing"
)

func TestCmd(t *testing.T) {

	//默认命令行操作
	cmd.SetDefault(cmd.Commend{
		Handler: func(args []string) {
			//这里直接调用内置 start 命令
			_ = cmd.Exec("test", args)
		},
	})

	//命令行示例
	cmd.Add(cmd.Commend{
		Name:    "test",
		Explain: "This is a test template.",
		Handler: func(args []string) {
			fmt.Println("Test command run successfully.")
		},
	})
}
