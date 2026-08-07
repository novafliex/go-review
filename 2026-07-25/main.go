package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := InitDB(); err != nil {
		log.Fatal("database initialization failed: ", err)
	}
	// 程序结束时关闭数据库连接池
	defer DB.Close()
	// 创建gin服务
	r := gin.Default()
	// 注册和登录路由
	auth := r.Group("/auth")
	{
		auth.POST("register", register)
		auth.POST("login", login)
	}
	// 启动服务
	if err := r.Run(":8081"); err != nil {
		log.Fatal("server start failed: ", err)
	}

}
