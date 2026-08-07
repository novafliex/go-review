package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	err := InitDB()
	if err != nil {
		fmt.Println("数据库初始化失败:", err)
		return
	}

	defer DB.Close()

	r := gin.Default()
	r.GET("/products", getProduct)
	r.GET("/products/:id", getProduct)
	r.POST("/products", createProduct)
	r.PUT("/products/:id", updateProduct)
	r.DELETE("/products/:id", deleteProduct)
	fmt.Println("启动服务器: http://localhost:8080")
	err = r.Run(":8080")
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
