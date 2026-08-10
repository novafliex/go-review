# 2026-07-26｜JWT 登录鉴权

## 当天练习

生成 JWT，用中间件保护资料接口并测试。

**和什么有关：** Claims、HS256、Bearer Token、中间件。

## 小项目结构

main.go 路由；handler.go 业务；middleware.go 鉴权；middleware_test.go 测试；utils.go Token。

## 快速复习与运行

```powershell
go run ./2026-07-26
go test ./2026-07-26
```

设置 JWT_SECRET 后运行和测试。

## 常见问题

密钥不能为空；限制签名算法；401 表示未认证。

