package cmd

import "errors"

type Commend struct {
	Name    string
	Explain string
	Handler func(args []string)
}

var cmdList = make(map[string]Commend)

// AddCmd 添加命令行
func AddCmd(cmd Commend) {
	cmdList[cmd.Name] = cmd
}

// ExecCmd 执行命令行（由框架完成此操作）
func ExecCmd(name string, args []string) error {
	if _, ok := cmdList[name]; !ok {
		return errors.New("no")
	}

	//执行匹配到的命令
	cmdList[name].Handler(args)
	return nil
}
