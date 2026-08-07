package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	err := InitDB()
	if err != nil {
		fmt.Println("数据库初始失败:", err)
		return
	}
	defer DB.Close()

	r := gin.Default()
	r.GET("/categories", getCategory)
	r.GET("/categories/:id", getCategory)
	r.PUT("/categories/:id", updateCategory)
	r.POST("/categories", createCategory)
	r.DELETE("/categories/:id", deleteCategory)
	//PATCH 是只修改资源的一部分
	r.PATCH("/categories/:id/status", updateCategoryStatus)
	fmt.Println("启动服务器:http://localhost:8084")
	err = r.Run(":8084")
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
