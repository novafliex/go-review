# Day 01：Go 并发

```text
2026-08-01/
├── examples/          # 3 个课程示例
├── exercises/         # 3 个练习
└── README.md          # 总结和练习要求
```

## 快速复习

- goroutine：并发执行任务。
- channel：在 goroutine 之间传递数据。
- WaitGroup：等待一组任务完成。
- Context：通知任务取消。
- Mutex/atomic：保护共享数据。

WaitGroup 只记这三句：

```text
Add(1)：多一个要等待的任务
Done()：一个任务完成
Wait()：等计数变成 0
```

注意：`go worker()` 只启动 goroutine，不会自动执行 `Add(1)`。

`examples/` 用于回顾课程代码，不再设置额外入口。

## 练习

### 1. 订单流水线（已完成）

```powershell
go test -race -count=1 -run TestOrderPipeline -v ./exercises
```

- 生产者关闭 `jobs`：它是唯一发送方。
- 协调者关闭 `results`：要先等所有 worker 退出。
- `Wait()` 放进另一个 goroutine：让主 goroutine 同时接收结果，避免死锁。
- `results` 关闭后，主 goroutine 的 `for range` 自动结束。

### 2. 安全计数器（已完成）

```powershell
go test -race -count=1 -run "TestMutexCounter|TestAtomicCounter" -v ./exercises
```

Mutex 和 atomic 版本都应稳定得到 100000。不安全版本只用于单独观察 race：

```powershell
$env:RUN_UNSAFE_RACE="1"
go test -race -count=1 -run TestUnsafeCounter -v ./exercises
Remove-Item Env:RUN_UNSAFE_RACE
```

### 3. Worker Pool（已完成）

```powershell
go test -race -count=1 -run TestWorkerPool -v ./exercises
```

固定数量 worker 处理任务，支持正常关闭和 Context 取消。
