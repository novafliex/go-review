package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// JSONResponse 自定义状态码、消息和数据
func JSONResponse(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// Success 默认返回 200
func Success(c *gin.Context, data any) {
	JSONResponse(
		c,
		http.StatusOK,
		"success",
		data,
	)
}

// ErrorResponse 返回错误信息
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	JSONResponse(
		c,
		statusCode,
		message,
		nil,
	)
}
