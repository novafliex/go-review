package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JSONResponse(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Response{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, data any) {
	JSONResponse(c, http.StatusOK, "success", data)
}

func Created(c *gin.Context, data any) {
	JSONResponse(c, http.StatusCreated, "created successfully", data)
}

func ErrorResponse(c *gin.Context, status int, message string) {
	JSONResponse(c, status, message, nil)
}
