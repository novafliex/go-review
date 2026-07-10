package main

import (
	"fmt"
	"net/http"
)

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("请求开始")
		fmt.Println("请求方法:", r.Method)
		fmt.Println("请求路径:", r.URL.Path)
		fmt.Println("请求参数:", r.URL.Query())
		fmt.Println("-------------------")
		next(w, r)
	}
}
