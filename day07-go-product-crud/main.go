package main

import (
	"fmt"
	"net/http"
)

var products = []Product{
	{"手机", 3000, 1},
	{"电脑", 34000.88, 2},
	{"相机", 4588, 3},
}

func main() {
	http.HandleFunc("/product", logMiddleware(product))

	fmt.Println("服务器启动成功!")
	fmt.Println("访问地址: http://localhost:8080/product")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
