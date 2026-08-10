# 2026-07-17｜音乐 REST CRUD

## 当天练习

用 Gorilla Mux 实现音乐内存 CRUD。

**和什么有关：** RESTful、JSON、资源建模。

## 小项目结构

main.go 包含模型、Handler 和路由。

## 快速复习与运行

```powershell
go run ./2026-07-17
```

服务地址 http://localhost:8080。

## 常见问题

创建通常返回 201；非法 JSON 会解码失败。

