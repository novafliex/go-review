# 2026-07-20｜Gin + MySQL 商品 CRUD

## 当天练习

用 database/sql 将商品持久化。

**和什么有关：** 连接池、环境变量、参数化 SQL。

## 小项目结构

main.go 入口；database.go 连接；handler.go SQL；model.go 模型；.env 配置。

## 快速复习与运行

```powershell
go run ./2026-07-20
```

先准备数据库和 .env。

## 常见问题

sql.Open 后还要 Ping；SQL 用 ? 占位；勿提交密码。

