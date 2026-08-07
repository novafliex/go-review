// package main 表示这是 Day 3 的最终完整 TCP 示例。
package main

import (
	// bufio 负责按行切分 TCP 字节流。
	"bufio"
	// context 负责传播停止信号。
	"context"
	// fmt 用于输出和按行回复。
	"fmt"
	// log 用于记录错误。
	"log"
	// net 提供 TCP 服务端能力。
	"net"
	// os 和 signal 用于监听 Ctrl+C。
	"os"
	"os/signal"
	// strings 负责命令解析。
	"strings"
	// sync 负责等待 worker、连接以及保存活动连接。
	"sync"
	// syscall 提供 SIGTERM。
	"syscall"
	// time 用于设置连接 deadline。
	"time"
)

const (
	// workerCount 限制同时执行消息业务的 goroutine 数量。
	workerCount = 3
	// maxLineBytes 限制客户端单条命令大小。
	maxLineBytes = 1024
	// readTimeout 限制客户端两条消息之间的最大等待时间。
	readTimeout = 30 * time.Second
	// writeTimeout 限制一次回复允许阻塞的时间。
	writeTimeout = 5 * time.Second
)

// messageJob 表示需要交给 worker 处理的一条命令。
type messageJob struct {
	// message 保存客户端发送的命令文本。
	message string
	// response 是当前连接等待处理结果的专属 channel。
	response chan string
}

// handleCommand 是 worker 真正执行的业务逻辑。
func handleCommand(message string) string {
	// 清理命令前后空白。
	message = strings.TrimSpace(message)
	// 根据协议返回结果。
	switch {
	case message == "PING":
		return "PONG"
	case strings.HasPrefix(message, "ECHO "):
		return strings.TrimPrefix(message, "ECHO ")
	default:
		return "ERR unknown command"
	}
}

// worker 不读取网络，只从 jobs 获取消息并执行命令解析。
func worker(id int, jobs <-chan messageJob, wg *sync.WaitGroup) {
	// worker 退出时通知 WaitGroup。
	defer wg.Done()
	// jobs 关闭后，for range 自动结束。
	for job := range jobs {
		// 处理当前命令。
		result := handleCommand(job.message)
		// 把结果发送回提交该任务的连接 goroutine。
		job.response <- result
		// 输出 worker 编号，方便观察任务如何分配。
		log.Printf("worker %d 处理命令 %q\n", id, job.message)
	}
}

// handleConnection 负责一个客户端的网络读写，并把业务交给 worker。
func handleConnection(ctx context.Context, conn net.Conn, jobs chan<- messageJob) {
	// 当前函数结束时关闭客户端连接。
	defer conn.Close()

	// 创建行协议 Scanner。
	scanner := bufio.NewScanner(conn)
	// 限制单行最大长度。
	scanner.Buffer(make([]byte, 256), maxLineBytes)

	// 持续处理当前连接上的消息。
	for {
		// 每次读取前刷新读 deadline。
		if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
			return
		}
		// 没读到完整行时结束连接。
		if !scanner.Scan() {
			return
		}

		// 为当前消息创建只能容纳一个结果的 channel。
		response := make(chan string, 1)
		// 组合消息和结果 channel，形成一个任务。
		job := messageJob{message: scanner.Text(), response: response}

		// 提交任务时同时监听服务停止信号。
		select {
		case jobs <- job:
			// worker 已经接收到任务。
		case <-ctx.Done():
			// 服务正在停止，不再提交新任务。
			return
		}

		// 等待 worker 结果，或者在停止时提前退出。
		select {
		case result := <-response:
			// 写回复前设置写 deadline。
			if err := conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
				return
			}
			// 把 worker 返回的结果写给客户端。
			if _, err := fmt.Fprintln(conn, result); err != nil {
				return
			}
		case <-ctx.Done():
			// 服务停止时不再等待处理结果。
			return
		}
	}
}

// main 组合监听器、连接 goroutine、worker pool 和优雅停止。
func main() {
	// 收到 Ctrl+C 或 SIGTERM 时取消 Context。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// main 退出时释放信号监听资源。
	defer stop()

	// 启动最终服务的监听器。
	listener, err := net.Listen("tcp", "127.0.0.1:8085")
	// 监听失败就结束程序。
	if err != nil {
		log.Fatal("启动 TCP 服务失败:", err)
	}
	// main 正常退出时也关闭监听器。
	defer listener.Close()

	// jobs 是连接 goroutine 向 worker 提交任务的通道。
	jobs := make(chan messageJob)
	// workerWG 用于等待所有 worker 退出。
	var workerWG sync.WaitGroup
	// 一次登记全部 worker。
	workerWG.Add(workerCount)
	// 启动固定数量的 worker。
	for id := 1; id <= workerCount; id++ {
		go worker(id, jobs, &workerWG)
	}

	// connections 保存所有活动 TCP 连接。
	var connections sync.Map
	// connectionWG 等待所有连接 goroutine。
	var connectionWG sync.WaitGroup

	// 监听停止信号并解除所有阻塞操作。
	go func() {
		// 等待 Context 被取消。
		<-ctx.Done()
		// 关闭 listener，使 Accept 返回。
		_ = listener.Close()
		// 关闭活动连接，使 Scanner 返回。
		connections.Range(func(key, _ any) bool {
			_ = key.(net.Conn).Close()
			return true
		})
	}()

	// 输出最终服务信息。
	fmt.Printf("最终 TCP 服务已启动：127.0.0.1:8085，worker=%d\n", workerCount)

	// 持续接收客户端，直到服务停止。
	for {
		// 等待一个连接。
		conn, err := listener.Accept()
		// 根据 Context 判断是否正常停止。
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Println("接收客户端失败:", err)
			continue
		}

		// 登记活动连接。
		connections.Store(conn, struct{}{})
		// 在启动 goroutine 前增加计数。
		connectionWG.Add(1)
		// 启动连接 goroutine。
		go func() {
			// 退出时减少连接计数。
			defer connectionWG.Done()
			// 退出时删除连接记录。
			defer connections.Delete(conn)
			// 开始处理当前客户端。
			handleConnection(ctx, conn, jobs)
		}()
	}

	// 先等待所有连接停止提交任务。
	connectionWG.Wait()
	// 此时不会再有发送方，安全关闭 jobs。
	close(jobs)
	// 等待 worker 处理完并退出。
	workerWG.Wait()
	// 最终确认所有并发任务已结束。
	fmt.Println("连接和 worker 全部退出，TCP 服务已停止")
}
