// package main 表示这是第三个可独立运行的 TCP 示例。
package main

import (
	// bufio 用换行符识别一条完整命令。
	"bufio"
	// fmt 用于打印和按行回复。
	"fmt"
	// log 用于记录错误。
	"log"
	// net 提供 TCP 能力。
	"net"
	// strings 用于清理字符串和判断命令前缀。
	"strings"
)

// handleCommand 只负责“输入命令 → 返回结果”，不处理网络。
func handleCommand(message string) string {
	// 去掉消息前后的空格、\r 和 \n。
	message = strings.TrimSpace(message)

	// 无表达式 switch 可以按条件从上到下判断。
	switch {
	// PING 常用来确认连接和服务是否正常。
	case message == "PING":
		return "PONG"
	// ECHO 后必须有一个空格和需要回显的文本。
	case strings.HasPrefix(message, "ECHO "):
		// 去掉包括空格在内的 ECHO 前缀，只返回后面的文本。
		return strings.TrimPrefix(message, "ECHO ")
	// 其他输入都不是当前协议支持的命令。
	default:
		return "ERR unknown command"
	}
}

// handleConnection 负责读取命令、调用解析函数并发送结果。
func handleConnection(conn net.Conn) {
	// 当前客户端处理结束时关闭连接。
	defer conn.Close()

	// 创建按行读取器，让 \n 成为应用层消息边界。
	scanner := bufio.NewScanner(conn)
	// 同一客户端可以连续发送多条命令。
	for scanner.Scan() {
		// 取得当前完整命令。
		message := scanner.Text()
		// 把命令交给不依赖网络的纯函数处理。
		response := handleCommand(message)
		// 打印输入和输出，方便学习时观察。
		fmt.Printf("客户端 %s：%q → %q\n", conn.RemoteAddr(), message, response)

		// 按行写回结果，让客户端也能用 ReadLine 接收。
		_, err := fmt.Fprintln(conn, response)
		// 写入失败时结束当前连接。
		if err != nil {
			log.Println("回复客户端失败:", err)
			return
		}
	}

	// 检查循环是否因读取错误结束。
	if err := scanner.Err(); err != nil {
		log.Println("读取客户端失败:", err)
	}
}

// main 负责启动命令协议服务。
func main() {
	// 监听本机 8082 端口。
	listener, err := net.Listen("tcp", "127.0.0.1:8082")
	// 监听失败时结束程序。
	if err != nil {
		log.Fatal("启动 TCP 服务失败:", err)
	}
	// main 退出时关闭监听器。
	defer listener.Close()

	// 输出支持的命令和地址。
	fmt.Println("命令协议服务已启动：127.0.0.1:8082，支持 PING 和 ECHO text")

	// 不断接受新连接。
	for {
		// Accept 等待一个客户端。
		conn, err := listener.Accept()
		// 当前接收失败就继续下一次循环。
		if err != nil {
			log.Println("接收客户端失败:", err)
			continue
		}

		// 每个连接由独立 goroutine 处理。
		go handleConnection(conn)
	}
}
