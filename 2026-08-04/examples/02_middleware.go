package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

type Middleware func(http.Handler) http.Handler

func apply(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func trace(name string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(name, "进入")
			next.ServeHTTP(w, r)
			fmt.Println(name, "退出")
		})
	}
}

func runMiddlewareExample() {
	final := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Println("Handler 执行业务")
		_, _ = w.Write([]byte("ok"))
	})
	handler := apply(final, trace("日志"), trace("认证"))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}
