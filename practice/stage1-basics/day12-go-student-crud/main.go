package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := InitDB(); err != nil {
		log.Printf("database initialization failed: %v", err)
		return
	}
	defer DB.Close()

	r := gin.Default()
	r.GET("/students", getStudents)
	r.GET("/students/:id", getStudent)
	r.POST("/students", createStudent)
	r.PUT("/students/:id", updateStudent)
	r.DELETE("/students/:id", deleteStudent)
	r.PATCH("/students/:id/status", updateStudentStatus)

	if err := r.Run(":8081"); err != nil {
		log.Printf("server startup failed: %v", err)
	}
}
