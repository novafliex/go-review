# go-review

这是一个按学习阶段整理的 Go 练习与项目仓库。基础练习、进阶练习和正式项目彼此分开，但统一使用根目录的 Go module：

```text
github.com/MikeJosion/go-review
```

## 目录结构

```text
go-review/
├── practice/
│   ├── stage1-basics/       # 已完成的基础与 Web 练习
│   └── stage2-advanced/     # 25 天计划 Day 1～4
├── projects/
│   ├── coursehub/           # 25 天计划 Day 5～24
│   └── coursehub-rebuild/   # Day 25 闭卷重建
├── docs/
│   ├── PLAN_25_DAYS.md
│   ├── LEARNING_STRUCTURE.md
│   └── DAILY_LOG_TEMPLATE.md
├── notes/                   # 每日复盘
├── go.mod
└── go.sum
```

## 第一阶段已有练习

| 目录 | 练习内容 |
|---|---|
| `day01-basic-syntax` | 变量、常量和数据类型 |
| `day02-functions` | 函数练习 |
| `day03-study-task` | 综合语法任务 |
| `day04-go-server` | `net/http` 与静态页面 |
| `day05-go-movie-crud` | 内存版电影 CRUD |
| `day06-go-music-crud` | 音乐 CRUD |
| `day07-go-product-crud` | 商品 CRUD 与中间件 |
| `day08-go-product-gin-crud` | Gin 商品 CRUD |
| `day09-go-product-mysql` | Gin + MySQL 商品服务 |
| `day10-go-product-user` | 商品、用户与前端联调 |
| `day11-go-category-crud` | 分类 CRUD |
| `day12-go-student-crud` | 学生 CRUD |
| `day13-go-bcrypt-basic` | bcrypt 密码摘要 |
| `day14-go-auth-bcrypt-nojwt` | 无 JWT 的登录鉴权 |
| `day15-go-jwt-auth-middleware` | JWT 鉴权中间件与测试 |

以上目录都位于 `practice/stage1-basics/`。

## 运行练习

在仓库根目录运行普通练习：

```powershell
go run ./practice/stage1-basics/day01-basic-syntax
go run ./practice/stage1-basics/day02-functions
go run ./practice/stage1-basics/day03-study-task
go run ./practice/stage1-basics/day15-go-jwt-auth-middleware
```

Day 04 使用相对静态资源路径，需要从自己的目录启动：

```powershell
Set-Location .\practice\stage1-basics\day04-go-server
go run .
```

数据库练习启动前需要准备本机 MySQL 和各练习需要的环境变量。真实密码只放在 `.env`，不要提交。

## 接下来怎么学

1. 打开 [25 天执行计划](./docs/PLAN_25_DAYS.md)。
2. 每天复制 [学习记录模板](./docs/DAILY_LOG_TEMPLATE.md) 为 `notes/dayNN.md`。
3. Day 1～4 在 `practice/stage2-advanced/` 写独立练习。
4. Day 5～24 始终在 `projects/coursehub/` 迭代同一个正式项目。
5. Day 25 在 `projects/coursehub-rebuild/` 闭卷重建。

更多整理原则见 [学习目录说明](./docs/LEARNING_STRUCTURE.md)。

