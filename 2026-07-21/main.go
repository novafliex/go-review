//只负责注册路由和启动服务器

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	err := InitDB()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	r.GET("/users", getUser)
	r.GET("/users/:id", getUser)
	r.GET("/products", getProduct)
	r.GET("/products/:id", getProduct)
	r.POST("/users", createUser)
	r.POST("/products", createProduct)
	r.PUT("/user/:id", updateUser)
	r.PUT("/products/:id", updateProduct)
	r.DELETE("/user/:id", deleteUser)
	r.DELETE("/products/:id", deleteProduct)
	fmt.Println("启动服务器: http://localhost:8080")
	r.Run(":8080")
}
