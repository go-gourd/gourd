
Gourd 1.1
===============

[![Go Report Card](https://goreportcard.com/badge/github.com/go-gourd/gourd)](https://goreportcard.com/report/github.com/go-gourd/gourd)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-gourd/gourd?status.svg)](https://pkg.go.dev/github.com/go-gourd/gourd?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/go-gourd/gourd/-/badge.svg)](https://sourcegraph.com/github.com/go-gourd/gourd?badge)
[![Release](https://img.shields.io/github/release/go-gourd/gourd.svg?style=flat-square)](https://github.com/go-gourd/gourd/releases)

```text
   _____                     _
  / ____|                   | |  Go       go1.21.1
 | |  __  ___  _   _ _ __ __| |  App      v1.1.0 (105)
 | | |_ |/ _ \| | | | '__/ _` |  Public   ./web/static
 | |__| | (_) | |_| | | | (_| |  Temp Dir ./runtime
  \_____|\___/ \__,_|_|  \__,_|  Log Dir  ./runtime/log
```
> Go-Gourd 是一个由事件驱动的快速开发框架，轻量、高效、易用。

## 主要特性
* 开箱即用
* 事件机制
* 定时任务
* 命令行
* 守护进程（仅linux）
* 日志工具

## 快速开始
使用git clone下载示例项目

```bash
git clone https://github.com/go-gourd/gourd-app.git
```

下载依赖项
```bash
go mod tidy
```

开始运行
```bash
go run main.go
```

运行成功，使用浏览器访问`http://localhsot:8888/`

## 参考文档

https://github.com/go-gourd/gourd/wiki

## 参与开发

直接提交PR或者Issue即可
## 版权信息

项目遵循Apache2开源协议发布，提供免费使用，项目包含的第三方源码和二进制文件之版权信息另行标注。

Email: `kyour@outlook.com` QQ: `2653907035`

更多细节参阅 [LICENSE](LICENSE)