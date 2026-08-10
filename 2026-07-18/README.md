# 2026-07-18｜原生 HTTP 商品 CRUD

## 当天练习

拆分模型、Handler、中间件和工具函数。

**和什么有关：** net/http、职责分层、中间件。

## 小项目结构

main.go 入口；handler.go 业务；middleware.go 日志；model.go 模型；utils.go 工具。

## 快速复习与运行

```powershell
go run ./2026-07-18
```

访问 http://localhost:8080/product。

## 常见问题

解析函数只解析参数，Handler 决定响应；错误后立即 return。

