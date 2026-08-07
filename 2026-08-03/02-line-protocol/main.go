// package main 表示当前目录可以用 go run 直接启动。
package main

import (
	// bufio.Scanner 可以按照换行符读取一条条消息。
	"bufio"
	// fmt 用于打印日志，也用于向连接写入一整行。
	"fmt"
	// log 用于记录网络错误。
	"log"
	// net 提供 TCP 监听和连接接口。
	"net"
)

// handleConnection 在同一条连接上持续处理多条消息。
func handleConnection(conn net.Conn) {
	// 当前客户端结束时关闭连接。
	defer conn.Close()

	// Scanner 默认以换行符 \n 作为一条消息的结束位置。
	scanner := bufio.NewScanner(conn)
	// Scan 每成功读取一整行，就进入一次循环。
	for scanner.Scan() {
		// Text 返回当前行的字符串，并去掉结尾的换行符。
		message := scanner.Text()
		// 在服务端终端观察收到的消息。
		fmt.Println("收到:", message)

		// Fprintln 写入回复，并自动在结尾补上换行符。
		_, err := fmt.Fprintln(conn, "服务器收到:", message)
		// 回复失败时无法继续使用当前连接。
		if err != nil {
			log.Println("回复客户端失败:", err)
			return
		}
	}

	// Scan 返回 false 后，检查是正常断开还是读取错误。
	if err := scanner.Err(); err != nil {
		// 有错误就记录，方便定位网络问题。
		log.Println("读取客户端失败:", err)
	}
}

// main 负责启动长连接行协议服务。
func main() {
	// 使用 8081，避免和第一个示例同时运行时端口冲突。
	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	// 监听失败就不能继续启动。
	if err != nil {
		log.Fatal("启动 TCP 服务失败:", err)
	}
	// 程序正常退出时关闭监听端口。
	defer listener.Close()

	// 输出启动信息。
	fmt.Println("行协议服务已启动：127.0.0.1:8081")

	// 持续接收不同客户端。
	for {
		// 等待一个客户端建立 TCP 连接。
		conn, err := listener.Accept()
		// 单次 Accept 失败不影响服务器继续工作。
		if err != nil {
			log.Println("接收客户端失败:", err)
			continue
		}

		// 每个连接单独处理，因此多个客户端不会互相等待。
		go handleConnection(conn)
	}
}
