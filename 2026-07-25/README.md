# 2026-07-25｜Gin 注册登录

## 当天练习

将 bcrypt 接入 Gin + MySQL。

**和什么有关：** 注册、登录、安全响应。

## 小项目结构

main.go 路由；database.go 数据库；handler.go 登录注册；utils.go 密码；model.go 模型；response.go 响应。

## 快速复习与运行

```powershell
go run ./2026-07-25
```

配置数据库后验证注册和登录。

## 常见问题

响应不能返回哈希；SQL 字段与 Scan 参数要对应。

