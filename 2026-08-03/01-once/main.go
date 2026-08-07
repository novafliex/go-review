// package main 表示这是一个可以直接运行的程序。
package main

import (
	// fmt 用于在终端输出信息。
	"fmt"
	// log 用于记录错误；Fatal 会打印错误并结束程序。
	"log"
	// net 提供 TCP 监听、连接、读取和写入能力。
	"net"
)

// handleConnection 处理一个客户端，并且只读取、回复一次。
func handleConnection(conn net.Conn) {
	// 函数结束时关闭当前客户端连接，避免连接资源泄漏。
	defer conn.Close()

	// 创建 1024 字节的缓冲区，用来接收网络传来的字节。
	buffer := make([]byte, 1024)
	// Read 把收到的数据写入 buffer；n 是本次实际读取的字节数。
	n, err := conn.Read(buffer)
	// 读取失败通常表示客户端断开或网络异常。
	if err != nil {
		// 记录错误后结束当前客户端的处理。
		log.Println("读取客户端失败:", err)
		return
	}

	// 只截取 buffer 中前 n 个有效字节，再转换为字符串。
	message := string(buffer[:n])
	// 在服务端终端打印客户端发送的内容。
	fmt.Println("收到:", message)

	// Write 只接受 []byte，所以先把回复字符串转换成字节切片。
	_, err = conn.Write([]byte("服务器收到"))
	// 写入失败说明回复没有成功发送给客户端。
	if err != nil {
		// 记录发送失败的原因。
		log.Println("回复客户端失败:", err)
	}
}

// main 是程序入口，只负责监听和分配客户端连接。
func main() {
	// 在本机 8080 端口启动 TCP 监听。
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	// 监听失败时服务器无法继续运行。
	if err != nil {
		// 打印错误并立即结束程序。
		log.Fatal("启动 TCP 服务失败:", err)
	}
	// main 正常退出前关闭监听器。
	defer listener.Close()

	// 告诉使用者服务已经启动，而不是程序卡住。
	fmt.Println("一次读写服务已启动：127.0.0.1:8080")

	// 无限循环让服务器可以不断接收新客户端。
	for {
		// Accept 会阻塞，直到有客户端连接进来。
		conn, err := listener.Accept()
		// 接收连接失败时记录错误，然后继续等待下一个连接。
		if err != nil {
			log.Println("接收客户端失败:", err)
			continue
		}

		// 为当前客户端启动独立 goroutine，主 goroutine 继续 Accept。
		go handleConnection(conn)
	}
}
