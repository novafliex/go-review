# 2026-07-22｜分类管理 CRUD

## 当天练习

分类增删改查与状态更新。

**和什么有关：** Gin、MySQL、唯一性、统一响应。

## 小项目结构

main.go 路由；database.go 数据库；category_handler.go 业务；model.go 模型；response.go 响应。

## 快速复习与运行

```powershell
go run ./2026-07-22
```

服务地址 http://localhost:8084。

## 常见问题

没查到同名分类表示可以创建；Scan 顺序要匹配 SQL。

