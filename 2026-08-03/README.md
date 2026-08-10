# 2026-08-03｜TCP 网络编程

## 当天练习

行协议、命令、保护、优雅停止和 worker pool。

**和什么有关：** TCP 字节流、边界、超时、Context。

## 小项目结构

01-once 到 06-worker-pool：六个递进小项目。

## 快速复习与运行

```powershell
go run ./2026-08-03/03-command-protocol
```

Listen → Accept → 读取 → 解析 → 回复。

## 常见问题

TCP 没有消息边界；Accept 和读取会阻塞；各项目端口不同。
