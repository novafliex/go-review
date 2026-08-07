// package main 表示这是一个可运行的优雅停止示例。
package main

import (
	// bufio 负责按行读取命令。
	"bufio"
	// context 传递“服务需要停止”的信号。
	"context"
	// fmt 用于输出和回复。
	"fmt"
	// log 用于记录错误。
	"log"
	// net 提供监听器和客户端连接。
	"net"
	// os 提供操作系统中断信号。
	"os"
	// signal 把 Ctrl+C 和 Context 取消关联起来。
	"os/signal"
	// strings 用于解析命令。
	"strings"
	// sync 提供 WaitGroup 和并发安全连接表。
	"sync"
	// syscall 提供程序终止信号 SIGTERM。
	"syscall"
)

// handleCommand 处理行协议命令。
func handleCommand(message string) string {
	// 清理命令两端空白。
	message = strings.TrimSpace(message)
	// 返回对应结果。
	switch {
	case message == "PING":
		return "PONG"
	case strings.HasPrefix(message, "ECHO "):
		return strings.TrimPrefix(message, "ECHO ")
	default:
		return "ERR unknown command"
	}
}

// handleConnection 持续处理一个客户端，连接关闭后自然退出。
func handleConnection(conn net.Conn) {
	// 函数结束时关闭连接。
	defer conn.Close()

	// 按行读取客户端命令。
	scanner := bufio.NewScanner(conn)
	// 每读到一行就处理一次。
	for scanner.Scan() {
		// 解析当前命令。
		response := handleCommand(scanner.Text())
		// 回复失败时结束当前连接。
		if _, err := fmt.Fprintln(conn, response); err != nil {
			return
		}
	}
}

// main 演示“停止接收 → 关闭旧连接 → 等待 goroutine”的顺序。
func main() {
	// Ctrl+C 或 SIGTERM 到来时，ctx 会自动变成已取消状态。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// main 退出时停止信号监听，释放内部资源。
	defer stop()

	// 启动 TCP 监听。
	listener, err := net.Listen("tcp", "127.0.0.1:8084")
	// 启动失败时直接结束。
	if err != nil {
		log.Fatal("启动 TCP 服务失败:", err)
	}
	// 正常路径下也保证监听器最终关闭。
	defer listener.Close()

	// connections 保存当前所有活动连接，停止时可以逐个关闭。
	var connections sync.Map
	// connectionWG 记录尚未退出的连接 goroutine 数量。
	var connectionWG sync.WaitGroup

	// 单独启动 goroutine 等待停止信号。
	go func() {
		// 没收到信号时一直阻塞在这里。
		<-ctx.Done()
		// 关闭 listener，让阻塞中的 Accept 立刻返回错误。
		_ = listener.Close()
		// 遍历并关闭活动连接，让阻塞中的 Scan 退出。
		connections.Range(func(key, _ any) bool {
			_ = key.(net.Conn).Close()
			return true
		})
	}()

	// 输出启动和停止方法。
	fmt.Println("优雅停止服务已启动：127.0.0.1:8084，按 Ctrl+C 停止")

	// 持续接收客户端，直到监听器因停止信号被关闭。
	for {
		// 等待新客户端。
		conn, err := listener.Accept()
		// Accept 返回错误时判断是不是正在停止。
		if err != nil {
			// Context 已取消说明这是正常停止，不需要继续 Accept。
			if ctx.Err() != nil {
				break
			}
			// 其他错误只记录，服务器继续运行。
			log.Println("接收客户端失败:", err)
			continue
		}

		// 把新连接登记进活动连接表。
		connections.Store(conn, struct{}{})
		// 在启动 goroutine 前把等待计数加一。
		connectionWG.Add(1)
		// 启动当前客户端处理 goroutine。
		go func() {
			// goroutine 退出时把等待计数减一。
			defer connectionWG.Done()
			// goroutine 退出时从活动连接表删除当前连接。
			defer connections.Delete(conn)
			// 执行真正的客户端读写。
			handleConnection(conn)
		}()
	}

	// 等待所有已经接收的客户端 goroutine 退出。
	connectionWG.Wait()
	// 到这里说明没有遗留连接任务。
	fmt.Println("所有连接已退出，TCP 服务已停止")
}
