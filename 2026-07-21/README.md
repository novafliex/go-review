# 2026-07-21｜商品与用户管理项目

## 当天练习

完成商品和用户 CRUD，并用单页前端联调。

**和什么有关：** Gin、MySQL、CORS、Vue 3 CDN。

## 小项目结构

main.go 路由；database.go 数据库；handler.go CRUD；model.go 模型；frontend/ 页面。

## 快速复习与运行

```powershell
go run ./2026-07-21
```

启动后端，再打开 frontend/index.html。

## 常见问题

确认 MySQL 和 .env；跨域需 CORS；数据库 ID 不是表格序号。

用户密码使用 bcrypt 保存，响应不返回密码；更新接口不允许修改角色。旧的明文测试密码需要重新设置。本例未实现登录鉴权，仅用于本地 CRUD 复习。
