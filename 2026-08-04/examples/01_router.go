package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func newGinRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	api.GET("/courses", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": []string{"Go 并发", "Go Web"}})
	})
	api.GET("/courses/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
	})
	return router
}

func runRouterExample() {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/courses/7", nil)
	response := httptest.NewRecorder()
	newGinRouter().ServeHTTP(response, request)
	fmt.Printf("路由参数响应: %s", response.Body.String())
}

func main() {
	fmt.Println("=== 1. HTTP 请求响应 ===")
	runHTTPExample()
	fmt.Println("\n=== 2. WebSocket 双向通信 ===")
	runWebSocketExample()
	fmt.Println("\n=== 3. Gin 路由 ===")
	runRouterExample()
	fmt.Println("\n=== 4. 中间件链 ===")
	runMiddlewareExample()
	fmt.Println("\n=== 5. RESTful API ===")
	runRESTExample()
}
