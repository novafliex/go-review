# 2026-08-04｜HTTP、WebSocket 与 API

## 当天练习

路由、中间件、REST、HTTP、WebSocket 和测试。

**和什么有关：** 请求响应、长连接、认证、统一错误。

## 小项目结构

examples/ 五个示例；exercises/ CRUD、Hub 与测试。

## 快速复习与运行

```powershell
go test -race ./2026-08-04/exercises
```

HTTP 适合请求响应，WebSocket 适合实时双向通信。

## 常见问题

401 是未认证，403 是权限不足；注意中间件顺序。
