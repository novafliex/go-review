# 2026-07-23｜学生管理 CRUD

## 当天练习

按 schema.sql 完成学生 CRUD 和状态更新。

**和什么有关：** Gin、MySQL、建表 SQL、errors.Is。

## 小项目结构

schema.sql 表；main.go 路由；database.go 数据库；student_handler.go 业务；model.go 模型；response.go 响应。

## 快速复习与运行

```powershell
go run ./2026-07-23
```

先执行 schema.sql。

## 常见问题

未建表会失败；用 errors.Is 判断 sql.ErrNoRows。

