package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var products = []Product{
	{"手机", 3000, 1},
	{"电脑", 34000.88, 2},
	{"相机", 4588, 3},
}

func main() {
	r := gin.Default()
	r.GET("/products", getProducts)
	r.GET("/products/:id", getProduct)
	r.POST("/products", createProduct)
	r.PUT("/products/:id", updateProduct)
	r.DELETE("/products/:id", deleteProduct)
	fmt.Println("启动服务器: http://localhost:8080")
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
