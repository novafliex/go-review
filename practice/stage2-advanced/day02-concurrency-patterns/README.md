# Day 02：并发模式

```text
day02-concurrency-patterns/
├── examples/    # 并发模式课程代码，3 个文件
├── exercises/   # 扇出、扇入、Pipeline 练习
└── README.md
```

## 快速复习

- 生产者-消费者：用 channel 协调生产和处理速度。
- 扇出：一个任务 channel，多个 goroutine 竞争处理。
- 扇入：多个生产者，结果汇总到一个 channel。
- Pipeline：任务拆成多个阶段，阶段之间用 channel 连接。

运行全部并发模式示例：

```powershell
go run ./examples
```

重点观察：

1. 每个 channel 由谁关闭？
2. 扇出中一个订单会被几个 worker 处理？
3. 扇入为什么要等所有生产者退出后再关闭输出？
4. Pipeline 的关闭信号怎样逐级传到最后？

运行练习测试：

```powershell
go test -race ./exercises
```
