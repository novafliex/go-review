package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run . 1 | 2 | 3 | practice")
		return
	}

	switch os.Args[1] {
	case "1":
		runLesson1()
	case "2":
		runLesson2()
	case "3":
		runLesson3()
	case "practice":
		runPractice()
	default:
		fmt.Println("只能输入: 1、2、3 或 practice")
	}
}
