
Gourd 0.3 Beta 2023
===============

Gourd is a Golang rapid development framework.

## 框架由来
我正在学习Golang，是一名PHP开发者，学习过程中发现Go在很多地方对于php开发习惯很不友好。  
为了从PHP更快上手Golang，我决定一边学习，一边构建一个自己用得惯的框架（代码结构）。  
也希望能够帮到正在学习Golang的同学。

## 主要特性
* 开箱即用
* 事件机制
* 定时任务
* 命令行
* 守护进程（开发中）

## 快速开始
使用git下载

```bash
git clone https://github.com/go-gourd/gourd-app.git
```

下载依赖
```bash
go mod tidy
```

开始运行
```bash
go run main.go
```

等待控制台输出
```text
   _____                     _ 
  / ____|                   | |  Go       go1.20
 | |  __  ___  _   _ _ __ __| |  Gourd    v0.2.0 (2)
 | | |_ |/ _ \| | | | '__/ _` |  Gin      v1.8.2
 | |__| | (_) | |_| | | | (_| |  Public   ./public
  \_____|\___/ \__,_|_|  \__,_|  Temp Dir ./runtime
--------------------------------------------------------
{"level":"info","time":"2023-02 *","msg":"Started http server. :8888"}
```
运行成功，使用浏览器访问`http://localhsot/`

## 参与开发

直接提交PR或者Issue即可

## 版权信息

Gourd遵循Apache2开源协议发布，并提供免费使用。

本项目包含的第三方源码和二进制文件之版权信息另行标注。

版权所有Copyright © 2020 by Kyour (http://www.kyour.cn) All rights reserved。

更多细节参阅 [LICENSE](LICENSE)