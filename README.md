# go-review

这是一个用于回顾和巩固 Go 语言基础的练习项目。

项目内容围绕 Go 的基础语法和后端开发常见场景展开，包含变量、数据类型、函数、结构体、切片、命令行交互、HTTP 服务、JSON 处理以及简单的 CRUD 接口实现。这个仓库不是一个完整的业务系统，而是我用来复习 Go 语法、整理练习代码、重新熟悉 Go 后端开发流程的学习记录。

## 项目目标

- 回顾 Go 基础语法和常用写法
- 熟悉函数、结构体、切片等核心概念
- 练习 Go 命令行程序的基本输入输出
- 复习 `net/http` 的基础使用
- 练习 REST API、JSON 编解码和简单 CRUD
- 为后续继续学习 Go 后端开发打基础

## 技术内容

项目中主要涉及：

- Go 基础语法
- Go Modules
- 标准库 `fmt`、`net/http`、`encoding/json`
- 第三方路由库 `github.com/gorilla/mux`
- 基于内存数据的简单 CRUD 示例

## 运行方式

进入项目目录后，可以按需运行不同练习：

```bash
go run ./day01-basic-syntax
go run ./day02-functions
go run ./day03-study-task
go run ./day04-go-server
go run ./day05-go-movie-crud
```

其中 HTTP 服务相关练习启动后，可以通过浏览器或 `curl` 访问对应接口。

`day06-go-music-crud` 是一个单独的 Go module，需要进入目录后运行：

```bash
cd day06-go-music-crud
go run .
```

## 当前状态

目前项目主要用于个人复习和练手，部分 CRUD 示例的数据保存在内存中，程序重启后会恢复初始数据。后续可以继续扩展数据库持久化、接口测试、配置管理和更完整的错误处理。
