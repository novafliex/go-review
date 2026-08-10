# 2026-08-01｜Go 并发基础

## 当天练习

goroutine、channel、WaitGroup、Context、Mutex、atomic、worker pool。

**和什么有关：** 并发协作、取消、竞态检测。

## 小项目结构

examples/ 完整示例；exercises/ 练习与测试。

## 快速复习与运行

```powershell
go test -race ./2026-08-01/exercises
```

Add 增加任务，Done 完成，Wait 等归零。

## 常见问题

发送方关闭 channel；共享变量需加锁或 atomic。
