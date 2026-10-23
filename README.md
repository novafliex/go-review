# Go 复习代码

记录日常 Go 复习中的示例、练习和笔记，按日期整理。

## 内容

- 基础语法、函数与综合练习
- HTTP、Gin、RESTful API 与 CRUD
- MySQL、bcrypt 与 JWT 鉴权
- goroutine、channel、并发模式与 TCP / WebSocket
- 数据库事务模板（Mock 模拟）

每个日期目录是一组独立练习，具体说明见目录内的 README。所有练习共用根目录的 `go.mod`。

## 运行

在仓库根目录执行，Go 版本要求见 `go.mod`。

```powershell
# 运行某一天的练习（同目录的 Go 文件需要一起编译）
go run ./2026-07-12

# 运行独立示例
go run ./2026-08-03/01-once

# 运行测试
go test ./...
```

`examples` 保存示例代码，部分需要通过同目录入口调用；`exercises` 中的练习通过测试验证。

MySQL 相关练习需要先准备数据库并配置 `MYSQL_DSN`，JWT 练习还需要设置 `JWT_SECRET`，具体见当天说明。服务示例建议逐个启动，避免端口冲突。
