# go-review 学习目录说明

唯一正确仓库：`E:\Code\Go\go-review`

## 整理后的结构

```text
go-review/
├── practice/
│   ├── stage1-basics/
│   │   ├── day01-basic-syntax/
│   │   ├── day02-functions/
│   │   ├── day03-study-task/
│   │   ├── ...
│   │   └── day15-go-jwt-auth-middleware/
│   └── stage2-advanced/
│       ├── day01-concurrency/
│       ├── day02-tcp/
│       ├── day03-web/
│       └── day04-rest/
├── projects/
│   ├── coursehub/
│   └── coursehub-rebuild/
├── docs/
├── notes/
├── go.mod
├── go.sum
└── README.md
```

## 分类规则

- `practice/stage1-basics/`：已经完成的基础语法、HTTP、CRUD、MySQL、bcrypt 和 JWT 练习。
- `practice/stage2-advanced/`：并发、TCP、HTTP/WebSocket 和中间件等单知识点进阶练习。
- `projects/coursehub/`：有完整业务目标的正式项目。项目按 `handler/service/repository` 和业务模块组织，不按日期拆散。
- `projects/coursehub-rebuild/`：Day 25 闭卷重建，开始前保持空实现。
- `docs/`：课程计划、结构约定和模板。
- `notes/`：每天自己的预测、错误、复盘和提交记录。

## Go module 规则

整个仓库继续共用根模块 `github.com/MikeJosion/go-review`。练习目录和 CourseHub 内都不要重复执行 `go mod init`，除非未来明确决定把 CourseHub 拆成独立仓库。

## 新增内容放哪里

- 只验证一个 Go 知识点：放 `practice/`。
- 会持续多天、有数据库和业务流程：放 `projects/`。
- 学习计划、架构说明、API 文档：放 `docs/` 或对应项目自己的 `docs/`。
- 当天心得和报错：放 `notes/`。

