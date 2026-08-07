// package main 表示这是一个可运行的“连接保护”示例。
package main

import (
	// bufio 用来按行读取，并限制单行最大长度。
	"bufio"
	// errors 用来判断 Scanner 错误中是否包含网络超时。
	"errors"
	// fmt 用于输出日志和回复客户端。
	"fmt"
	// log 用于记录错误。
	"log"
	// net 提供 TCP 连接以及 net.Error 超时接口。
	"net"
	// strings 用于解析命令。
	"strings"
	// time 用于计算读写截止时间。
	"time"
)

const (
	// maxLineBytes 防止客户端发送超大单行占满服务器内存。
	maxLineBytes = 1024
	// readTimeout 防止客户端连接后一直不发送数据。
	readTimeout = 30 * time.Second
	// writeTimeout 防止客户端不读取回复导致服务器一直阻塞。
	writeTimeout = 5 * time.Second
)

// handleCommand 解析当前行协议支持的命令。
func handleCommand(message string) string {
	// 清理两端空白字符。
	message = strings.TrimSpace(message)
	// 按命令内容返回结果。
	switch {
	case message == "PING":
		return "PONG"
	case strings.HasPrefix(message, "ECHO "):
		return strings.TrimPrefix(message, "ECHO ")
	default:
		return "ERR unknown command"
	}
}

// handleConnection 给连接加入长度限制和读写超时。
func handleConnection(conn net.Conn) {
	// 当前客户端结束时释放连接。
	defer conn.Close()

	// 创建按行读取器。
	scanner := bufio.NewScanner(conn)
	// 初始缓冲区较小，第二个参数规定允许的最大 token 长度。
	scanner.Buffer(make([]byte, 256), maxLineBytes)

	// 持续处理当前客户端的每一条命令。
	for {
		// 每次读取前重新设置 deadline，活跃客户端可以持续使用连接。
		if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
			log.Println("设置读取超时失败:", err)
			return
		}

		// Scan 等到完整一行、连接断开、超时或行长度超限。
		if !scanner.Scan() {
			break
		}

		// 解析当前命令并得到回复。
		response := handleCommand(scanner.Text())
		// 每次写入前设置写 deadline。
		if err := conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
			log.Println("设置写入超时失败:", err)
			return
		}
		// 把结果作为完整一行发送给客户端。
		if _, err := fmt.Fprintln(conn, response); err != nil {
			log.Println("回复客户端失败:", err)
			return
		}
	}

	// 读取循环结束后获取具体原因。
	if err := scanner.Err(); err != nil {
		// net.Error 可以告诉我们错误是否由 deadline 超时引起。
		var netErr net.Error
		// 超时与非法超长输入使用不同日志，便于排查。
		if errors.As(err, &netErr) && netErr.Timeout() {
			log.Println("客户端读取超时:", conn.RemoteAddr())
		} else {
			log.Println("消息过长或读取失败:", err)
		}
	}
}

// main 启动带连接保护的 TCP 服务。
func main() {
	// 监听本机 8083 端口。
	listener, err := net.Listen("tcp", "127.0.0.1:8083")
	// 监听失败时不能继续运行。
	if err != nil {
		log.Fatal("启动 TCP 服务失败:", err)
	}
	// main 退出时关闭监听器。
	defer listener.Close()

	// 输出本示例启用的保护参数。
	fmt.Printf("保护版服务已启动：127.0.0.1:8083，最大行=%d，读超时=%s\n", maxLineBytes, readTimeout)

	// 持续接收客户端。
	for {
		// 等待新连接。
		conn, err := listener.Accept()
		// 单次失败时记录并继续。
		if err != nil {
			log.Println("接收客户端失败:", err)
			continue
		}

		// 每个连接独立执行读写与超时控制。
		go handleConnection(conn)
	}
}
