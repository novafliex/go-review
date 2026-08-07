# Day 03：TCP

今天只学习 TCP。代码按难度逐步增加，每个类别只有一个可独立运行的 `main.go`。

```text
2026-08-03/
├── 01-once/main.go             # 一次读取、一次回复
├── 02-line-protocol/main.go    # 长连接与换行协议
├── 03-command-protocol/main.go # PING、ECHO、未知命令
├── 04-protection/main.go       # 最大行长度、读写超时
├── 05-graceful-stop/main.go    # Context、信号、WaitGroup
└── 06-worker-pool/main.go      # 固定 worker 数处理消息
```

## 学习顺序

每一类都按相同方式学习：

```text
看注释理解 → 运行验证 → 只写流程 → 关掉答案复现新增部分
```

不要每次重写全部代码。例如学到第 3 类时，只复现 `handleCommand()`；学到第 5 类时，只复现优雅停止相关代码。

## 运行

从仓库根目录选择一个类别运行：

```powershell
go run ./2026-08-03/01-once
go run ./2026-08-03/02-line-protocol
go run ./2026-08-03/03-command-protocol
go run ./2026-08-03/04-protection
go run ./2026-08-03/05-graceful-stop
go run ./2026-08-03/06-worker-pool
```

六个服务分别使用 `8080`～`8085` 端口，运行后按 `Ctrl+C` 停止。

## 验证客户端

把 `$port` 改成当前服务使用的端口：

```powershell
$port = 8082
$client = [System.Net.Sockets.TcpClient]::new("127.0.0.1", $port)
$stream = $client.GetStream()
$utf8 = [System.Text.UTF8Encoding]::new($false)
$writer = [System.IO.StreamWriter]::new($stream, $utf8, 1024, $true)
$reader = [System.IO.StreamReader]::new($stream, $utf8, $false, 1024, $true)
$writer.AutoFlush = $true
$writer.WriteLine("PING")
$reader.ReadLine()
$writer.WriteLine("ECHO hello")
$reader.ReadLine()
$client.Close()
```

## 最终要能说出的流程

```text
Listen → Accept → goroutine → 按行读取 → 解析命令 → worker 处理 → 按行回复
                                      ↓
                         超时、长度限制、Context 停止
```

- TCP 是字节流，没有消息边界；本练习用换行符定义边界。
- `Accept()` 等连接，`Scan()` 等一行数据，它们都会阻塞。
- 每个客户端由一个 goroutine 负责读写，消息计算可以交给固定 worker pool。
- 生产者关闭自己唯一发送的 channel；服务停止时先不接新连接，再等旧连接退出。
