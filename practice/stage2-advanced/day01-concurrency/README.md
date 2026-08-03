# Day 01：并发

这里只有 5 个 Go 文件：入口、三节课、一个练习。

```text
day01-concurrency/
├── main.go
├── lesson01_concurrency.go
├── lesson02_channel.go
├── lesson03_sync.go
├── practice.go
└── README.md
```

在当前目录运行：

```powershell
go run . 1          # 第一节：并发、限流、Context
go run . 2          # 第二节：Channel、订单流水线、select
go run . 3          # 第三节：RWMutex、atomic
go run . practice   # 查看今天的独立练习
```

写完练习后检查：

```powershell
go run -race . practice
```
