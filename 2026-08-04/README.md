# Day 04：HTTP、WebSocket 与 Web API

Day 3 已经只保留 TCP。第二阶段剩余第 6～11 节统一放在这里，学完后进入第三阶段项目实战。

```text
2026-08-04/
├── examples/
│   ├── 01_router.go       # Gin 路由，也是示例运行入口
│   ├── 02_middleware.go   # 中间件链
│   ├── 03_restful.go      # RESTful 请求与响应
│   ├── 04_http.go         # net/http
│   └── 05_websocket.go    # WebSocket
├── exercises/             # HTTP、Hub、CRUD、中间件及测试
└── README.md
```

## 快速复习

- HTTP：一次请求对应一次响应，Handler 读取请求并写回状态码、响应头和响应体。
- WebSocket：通过 HTTP 升级成长连接，客户端和服务端都可以主动发送消息。
- Web 开发流程：需求 → 设计 → 开发 → 测试与上线。
- 路由：根据 HTTP 方法和路径，把请求交给对应 Handler。
- 中间件：接收一个 Handler，返回增强后的新 Handler。
- RESTful：路径表示资源，HTTP 方法表示操作，状态码表示结果。

## 运行

```powershell
go run ./2026-08-04/examples
go test -race ./2026-08-04/exercises
```

## CRUD 路由

| 方法 | 路径 | 含义 | Bearer Token |
|---|---|---|---|
| GET | `/courses` | 课程列表 | 不需要 |
| POST | `/courses` | 创建课程 | 需要 |
| GET | `/courses/{id}` | 课程详情 | 不需要 |
| PUT | `/courses/{id}` | 更新课程 | 需要 |
| DELETE | `/courses/{id}` | 删除课程 | 需要 |

中间件顺序：

```text
Request ID → Access Log → Recover → Auth → Handler
```

统一错误格式：

```json
{"code":"NOT_FOUND","message":"course not found","request_id":"req-1"}
```

## 复习问题

1. HTTP 与 WebSocket 的边界是什么？HTTP 适合请求响应，WebSocket 适合实时双向推送。
2. `401` 和 `403` 有什么不同？`401` 是未认证，`403` 是已认证但权限不足。
3. 为什么中间件倒序组合？这样参数中的第一个中间件才能成为最外层。
4. 创建资源为什么返回 `201`？它明确表示服务器成功创建了新资源。
