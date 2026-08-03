# 25 天 Go 进阶 + CourseHub 企业项目实战计划

## 执行假设

本计划按每天 **3～4 小时**设计：课程约 22 小时，编码、测试和复盘约 55～65 小时。八套课程的所有分集会在 Day 24 前看完，Day 25 完成闭卷重建和总验收。

每天固定节奏：

1. **10～20 分钟，先写后看**：只看当天标题，先写接口、伪代码或自己的实现猜想。
2. **45～75 分钟，完整听课**：推荐 1.0～1.25 倍速，所有指定分集都看；遇到 coding 段落暂停预测下一步。
3. **90～150 分钟，关掉视频重写**：允许查官方文档和报错，不允许照抄讲师成品。
4. **20～40 分钟，测试和复盘**：主动测一个失败场景，填写日志并提交 Git。

正确仓库已经有根模块 `github.com/MikeJosion/go-review`，所有练习和 CourseHub 共用它，不重复执行 `go mod init`。当天验收的通用底线在仓库根目录执行：

```powershell
gofmt -w .
go test ./...
go vet ./...
git status --short
```

涉及并发的当天再执行：

```powershell
go test -race ./...
```

> 只看完视频不算完成。代码能跑但没有失败场景测试，也不算真正掌握。

## 课程清单与总量

1. [第二阶段：4 小时 Go 语言进阶](https://www.bilibili.com/video/BV16bsGzEEDn/)：11 节，约 3 小时 48 分。
2. [第三阶段：需求分析、技术调研、技术方案、项目管理](https://www.bilibili.com/video/BV1kb1vBJEi2/)：14 节，约 3 小时 11 分。
3. [第三阶段：项目代码框架、鉴权、日志、中间件](https://www.bilibili.com/video/BV1DuCCB7ECu/)：9 节，约 2 小时 17 分。
4. [第三阶段：手机号、飞书登录](https://www.bilibili.com/video/BV128UjB2ESQ/)：14 节，约 4 小时 27 分（原列表编号跳过 11，以页面实际分集为准）。
5. [第三阶段：管理员、角色、权限与 RBAC](https://www.bilibili.com/video/BV1xuUfBwE9E/)：6 节，约 1 小时 45 分。
6. [第三阶段：录播目录和对象存储](https://www.bilibili.com/video/BV1vc2EBtETn/)：6 节，约 1 小时 54 分。
7. [第三阶段：录播课时管理](https://www.bilibili.com/video/BV18qmtBmEwr/)：6 节，约 1 小时 47 分。
8. [第三阶段：课程管理](https://www.bilibili.com/video/BV1jx6tB6EMR/)：9 节，约 2 小时 56 分。

合计约 **22 小时 4 分**。

## 最终项目：CourseHub

做一个企业内部课程管理 API，围绕同一条业务链迭代，而不是每天换新项目。

核心能力：

- 用户支持密码、手机验证码和飞书三种登录入口；外部服务通过接口隔离，本地可用 fake adapter 验证完整流程。
- JWT/会话鉴权、登出和重置密码。
- 用户—角色—权限—菜单的 RBAC。
- 课程、树形目录、录播课时的增删改查、分页和事务。
- 对象存储抽象，本地文件实现必须可运行；腾讯云 COS 实现由环境变量启用。
- 结构化日志、request ID、恢复、鉴权、RBAC 等中间件。
- MySQL + GORM，核心逻辑有单元测试和集成测试。

仓库已经按“练习”和“项目”分区。Day 5 后建议在 `projects/coursehub/` 内逐步形成以下结构（不要第一天一次性建完）：

```text
projects/coursehub/
├── cmd/api/
├── internal/
│   ├── config/       ├── platform/database/
│   ├── logger/       ├── middleware/
│   ├── auth/         ├── user/
│   ├── rbac/         ├── course/
│   ├── catalog/      ├── lesson/
│   └── storage/
├── docs/
├── migrations/
└── scripts/
```

## 第一阶段：进阶基础（Day 1～4）

### Day 1：并发模型、channel、同步原语

**听课（约 72 分钟）**

- 第二阶段 1《并发模型》
- 第二阶段 2《channel 管道》
- 第二阶段 3《同步原语》

**亲手敲**

- 原有 `day01-basic-syntax`～`day15-go-jwt-auth-middleware` 已归档到 `practice/stage1-basics/`；新建并进入 `practice/stage2-advanced/day01-concurrency/`，继续使用仓库根 `go.mod`。
- `day01-concurrency/pool.go`：实现有界 worker pool，包含 `jobs`、`results`、`context.Context` 取消和正常关闭。
- `day01-concurrency/counter.go`：分别用 `sync.Mutex`、`atomic` 实现安全计数器，并故意保留一个仅供测试演示的非安全版本。
- `day01-concurrency/pool_test.go`：验证任务不丢失、取消后退出、重复关闭不 panic；用 `-race` 看见并修复数据竞争。

**验收**

- 能口述 goroutine 泄漏的两个原因、channel 由谁关闭、Mutex 与 channel 的选择边界。
- 在仓库根目录执行 `go test -race ./practice/stage2-advanced/day01-concurrency/...` 通过。
- Git：`day01: implement cancellable worker pool`。

### Day 2：并发模式与 TCP/IP

**听课（约 48 分钟）**

- 第二阶段 4《并发模式》
- 第二阶段 5《TCP-IP 网络基础》

**亲手敲**

- `practice/stage2-advanced/day02-tcp/server.go`：实现行协议 TCP 服务，支持 `PING`、`ECHO text`、未知命令错误。
- 每连接一个 goroutine；加入读写 deadline、最大行长度、优雅停止，并复用 Day 1 的 worker pool 处理消息。
- `day02-tcp/server_test.go`：使用 `net.Pipe` 或随机本地端口测试半包/多包、超时、客户端断开和服务停止。

**验收**

- 能画出客户端连接到响应的完整时序，并说明 TCP 为什么没有“消息边界”。
- 在仓库根目录执行 `go test -race ./practice/stage2-advanced/day02-tcp/...` 通过。
- Git：`day02: build concurrent tcp line server`。

### Day 3：HTTP、WebSocket 与 Web 产品流程

**听课（约 50 分钟）**

- 第二阶段 6《http 编程》
- 第二阶段 7《websocket 初探》
- 第二阶段 8《Web 服务产品开发流程》

**亲手敲**

- `practice/stage2-advanced/day03-web/server.go`：只用 `net/http` 写 `/healthz`、`GET /courses`、`POST /courses`。
- `day03-web/hub.go`：WebSocket 连接中心；创建课程后向所有客户端广播事件；断连必须注销，避免 goroutine 泄漏。
- 先写 1 页 `day03-web/mini-product-flow.md`：用户故事、接口、数据、失败场景、验收标准，再编码。
- 用 `httptest` 测状态码、JSON、非法请求体、方法错误；为 hub 写连接/断开测试。

**验收**

- 能解释 HTTP 与 WebSocket 的边界，以及 handler 为什么不应直接承载全部业务逻辑。
- Git：`day03: add http api and websocket events`。

### Day 4：路由、中间件、RESTful

**听课（约 58 分钟）**

- 第二阶段 9《路由系统》
- 第二阶段 10《中间件》
- 第二阶段 11《RESTful》

**亲手敲**

- 为 Day 3 API 加完整的内存 CRUD：列表、详情、创建、更新、删除。
- 自己实现中间件链：request ID、耗时日志、recover、简单 Bearer token。
- 统一 JSON 错误结构：`code`、`message`、`request_id`；正确使用 400/401/404/409/500。
- 表驱动测试验证中间件执行顺序、panic 恢复、未授权与资源不存在。

**验收**

- 不看代码写出一个中间件签名和 REST 资源路由表。
- `practice/stage2-advanced/` 全部测试及 race test 通过。
- Git：`day04: complete restful api middleware chain`。

## 第二阶段：从需求到可运行框架（Day 5～10）

### Day 5：需求、技术文档、MySQL、Git、配置

**听课（约 69 分钟）**

- 全流程 1～5：《产品需求分析》《技术文档框架》《Mysql 介绍和简单使用》《git 协作流程简单使用》《服务配置和环境变量》

**亲手敲**

- 新建 `projects/coursehub/` 的业务目录；从此 Day 5～24 都在这个正式项目中工作，但继续复用仓库根 `go.mod`，不要再初始化嵌套 module。
- `docs/requirements.md`：角色、6 条用户故事、范围/非范围、验收标准。
- `docs/domain.md`：User、Role、Permission、Course、Catalog、Lesson、Object 的关系和不变量。
- `internal/config/config.go`：只从环境变量读取配置，校验端口、DSN、JWT key；写 `.env.example`，真实 `.env` 已由仓库根 `.gitignore` 忽略。
- MySQL 手工创建 `coursehub` 和单独的开发用户；不要把 root 密码写进仓库。
- 为配置写表驱动测试：缺变量、非法值、正常值。

**验收**

- 每个需求都能对应到至少一个可测试验收标准。
- Git：`day05: define coursehub requirements and config`。

### Day 6：GORM/gentool、接口管理与技术方案

**听课（约 89 分钟）**

- 全流程 6～11：《gorm 和 gentool》《ApiFox 和接口管理》《技术架构设计》《模块设计》《数据库结构设计》《接口设计》

**亲手敲**

- `docs/architecture.md`：画依赖方向 `handler -> service -> repository`，说明为何外部服务要 adapter。
- `docs/schema.sql` 或 `migrations/001_init.sql`：先设计 users、roles、permissions、user_roles、role_permissions、courses、catalogs、lessons、objects。
- 建 `internal/platform/database/`，用 GORM 建立 MySQL 连接、连接池和关闭逻辑；gentool 只做一次实验，最终生成物必须读懂。
- 在 Apifox（或仓库内 OpenAPI 文件）定义第一版接口；为 repository 写最小的 create/get 集成测试。

**验收**

- 能说明事务边界、唯一索引、外键/软删除的取舍。
- Git：`day06: design schema and gorm repository`。

### Day 7：后端服务骨架与项目计划

**听课（约 34 分钟）**

- 全流程 12～14：《后端服务框架》《技术方案评审》《项目管理和项目计划》

**亲手敲**

- 引入课程使用的 Web 框架（若课程是 Gin，就用 Gin），建立 `cmd/api/main.go`。
- 完成配置 -> 日志 -> 数据库 -> router -> server 的依赖装配，加入 `/healthz`、`/readyz`。
- `docs/plan.md` 记录剩余 18 天里程碑、风险和完成定义；`docs/review-checklist.md` 写评审问题。
- 处理 SIGINT/SIGTERM，服务必须在超时内优雅退出。

**验收**

- 本地能启动、访问健康检查并正常停止；错误 DSN 时给出可定位的错误。
- Git：`day07: bootstrap coursehub api service`。

### Day 8：代码框架一～三

**听课（约 55 分钟）**

- 代码框架 1～3：《项目代码框架一》《代码框架二》《代码框架三 new》

**亲手敲**

- 完成 `internal/user` 的第一条垂直切片：model -> repository -> service -> handler -> route。
- 实现创建用户与按 ID 查询；定义领域错误，不让 handler 判断数据库错误文本。
- 使用构造函数注入依赖；service 单测使用手写 fake repository。

**验收**

- 能追踪一个请求穿过每层的输入/输出，并说出每层不该做什么。
- Git：`day08: implement first user vertical slice`。

### Day 9：代码框架四～六

**听课（约 39 分钟）**

- 代码框架 4～6：《项目代码框架四》《项目代码框架五》《项目代码框架六》

**亲手敲**

- 加统一响应、参数校验、错误映射和分页类型。
- 把 Day 8 的用户功能补成列表、更新、禁用；数据库更新只允许白名单字段。
- handler 用 `httptest`，repository 用测试数据库，service 用 fake；三层各测不同责任。

**验收**

- 非法 JSON、重复手机号、缺失用户、数据库失败都有明确测试。
- Git：`day09: harden user layers and error mapping`。

### Day 10：日志与鉴权中间件

**听课（约 43 分钟）**

- 代码框架 7～9：《日志中间件》《日志中间件 1》《鉴权中间件》

**亲手敲**

- `internal/logger/`：结构化日志，字段至少含 request_id、method、path、status、latency、user_id。
- `internal/middleware/`：request ID、access log、recover、JWT auth；日志中禁止出现密码、验证码和完整 token。
- JWT 校验 issuer、expiry、签名算法；把 claims 安全放进请求上下文。
- 中间件测试覆盖：无 token、过期、错误签名、合法 token、handler panic。

**验收**

- 能从一次失败请求的 request ID 找到完整日志链路。
- Git：`day10: add structured logging and jwt middleware`。

## 第三阶段：登录与 RBAC（Day 11～16）

### Day 11：登录需求与滑块校验

**听课（约 73 分钟）**

- 登录 1～4：《功能需求实现》《滑块校验 1》《滑块校验 2》《滑块校验 3》

**亲手敲**

- `internal/auth/` 定义 `CaptchaVerifier`、`CodeSender`、`OAuthProvider` 接口。
- 实现可运行的本地 fake 滑块校验，不依赖第三方密钥；校验票据一次性使用并有过期时间。
- 写登录状态机/流程图，明确校验、查询用户、签发会话的顺序。
- 测试：过期票据、重复使用、伪造票据、下游失败。

**验收**

- 更换校验供应商时 service 不需要修改。
- Git：`day11: model auth flow and captcha adapter`。

### Day 12：手机号与密码登录

**听课（约 59 分钟）**

- 登录 5～7：《用户手机号登录 1》《用户手机号登录 2》《用户手机号登录 3》

**亲手敲**

- 用户表加入唯一手机号、密码摘要、状态；密码使用 bcrypt/argon2，不存明文。
- 实现 `POST /auth/login/password` 和签发 access token；错误信息不能泄露“手机号是否存在”。
- 对同一手机号/IP 加登录失败计数与短时限流；成功后清除计数。
- 测试正确、错密码、禁用用户、重复手机号、暴力尝试。

**验收**

- 能解释哈希与加密的区别，以及为何登录错误提示要模糊化。
- Git：`day12: implement secure phone password login`。

### Day 13：验证码与飞书登录

**听课（约 69 分钟）**

- 登录 8～10：《手机号验证码登录 1》《手机号验证码登录 2—飞书扫码》《手机号验证码登录 3—飞书验证码》

**亲手敲**

- `POST /auth/code/send`、`POST /auth/login/code`：验证码摘要存储、TTL、次数限制、成功即失效。
- 实现 fake 短信 sender；实际短信服务只能作为环境变量启用的 adapter。
- 实现 Feishu OAuth adapter 的接口和本地 fake 回调，校验 state 防 CSRF；真实账号凭证可选，不影响本地验收。
- 测试重放、过期、次数超限、错误 state、供应商超时。

**验收**

- 关闭第三方网络，完整验证码/飞书 fake 流程仍可自动测试。
- Git：`day13: add code and feishu login adapters`。

### Day 14：登录收尾、绑定、登出、重置密码

**听课（约 68 分钟）**

- 登录页面剩余全部分集：《手机号验证码登录 4—完成》《多种登录方式—菜单可配置》《飞书扫码绑定和解绑》《登出系统和重置密码》

**亲手敲**

- 登录方式由配置决定，但后端仍强制执行安全规则。
- 建 account_bindings/session 或 refresh token 数据模型；实现飞书绑定/解绑、登出撤销、重置密码后撤销旧会话。
- 对绑定、解绑、重置密码写越权测试和幂等测试。
- 更新 OpenAPI/Apifox 登录集合，跑通三种登录演示。

**验收**

- 被撤销 token 不可再访问；A 用户不能解绑 B 用户账号。
- Git：`day14: finish account lifecycle and token revocation`。

### Day 15：协程池与菜单权限

**听课（约 55 分钟）**

- RBAC《管理员角色权限—协程池》《菜单权限配置一》《菜单权限配置二》

**亲手敲**

- 把 Day 1 worker pool 变成审计事件处理器，服务退出前完成或明确丢弃策略。
- 实现菜单树 CRUD、父子关系校验、排序；禁止循环父子关系。
- 权限常量采用 `course:read`、`course:write` 等资源:动作形式。
- 测试树构造、孤儿节点、循环、并发写审计队列。

**验收**

- race test 通过；能说明异步审计的可靠性取舍。
- Git：`day15: implement menu tree and audit workers`。

### Day 16：角色管理和 RBAC 授权

**听课（约 50 分钟）**

- RBAC 剩余全部：《角色管理一》《角色管理二》《给角色分配权限》

**亲手敲**

- 实现角色 CRUD、用户分配角色、角色分配权限，分配操作使用事务。
- `RequirePermission("course:write")` 中间件；每次请求从可信身份读取 user ID。
- 写权限矩阵测试：游客、普通教师、管理员；至少包含水平越权和垂直越权。
- 变更角色后旧权限如何生效要形成明确策略（即时查库或短 TTL 缓存）。

**验收**

- 权限拒绝返回 403，未登录返回 401，两者不混用。
- Git：`day16: enforce transactional rbac permissions`。

## 第四阶段：目录、对象存储、课时、课程（Day 17～24）

### Day 17：录播目录与对象存储概念

**听课（约 50 分钟）**

- 目录/对象存储 1～3：《录播目录管理一》《录播目录管理二》《腾讯云对象存储讲解》

**亲手敲**

- `internal/catalog/`：课程下的树形目录 CRUD，校验课程归属、唯一顺序和最大深度。
- `internal/storage/` 定义 `Put/Get/Delete/Presign` 接口和 Object 元数据。
- 先实现本地文件 adapter，路径必须清洗，阻止 `../` 目录穿越。
- 测目录越权、循环/过深、本地上传下载删除、非法路径。

**验收**

- 业务 service 不依赖腾讯云 SDK；本地 adapter 能跑完整测试。
- Git：`day17: add catalog and storage abstraction`。

### Day 18：COS adapter 与课时创建

**听课（约 65 分钟）**

- 目录/对象存储 4～6：《对象存储骨架代码》《腾讯云对象存储实现》《完成录播课时创建》

**亲手敲**

- 实现腾讯云 COS adapter 骨架；没有个人凭证时用 fake server 合同测试，绝不把密钥写入代码。
- 实现预签名上传流程：申请上传 -> 客户端上传 -> 回调/确认 -> 创建 lesson。
- 校验 object 所有权、类型、大小和状态；数据库失败时设计对象补偿清理任务。
- 测伪造 object key、重复确认、数据库失败、对象服务超时。

**验收**

- 能解释为什么不让大视频经过 API 服务中转，以及 DB 与对象存储为何不能做普通数据库事务。
- Git：`day18: implement object upload and lesson creation`。

### Day 19：录播课时列表

**听课（约 46 分钟）**

- 课时管理 1～2：《录播课时列表一》《录播课时列表二》

**亲手敲**

- `GET /courses/{courseID}/lessons`：分页、目录过滤、标题搜索、排序。
- repository 查询避免 N+1；返回 DTO，不直接暴露 GORM model。
- 测 page 边界、空结果、跨课程目录、稳定排序；记录一条查询 SQL 并分析索引。

**验收**

- 同样数据同样排序参数必须稳定返回；跨课程数据不可见。
- Git：`day19: add paginated lesson queries`。

### Day 20：课时编辑、修复与视频清理

**听课（约 62 分钟）**

- 课时管理 3～6：《录播课时编辑一》《录播课时编辑二》《录播课时 bug 修复》《录播课时视频更新清理》

**亲手敲**

- 实现课时更新；替换视频时事务内更新引用，事务成功后异步清理旧对象。
- 删除失败进入可重试 cleanup_jobs，任务必须幂等并有限退避。
- 给课程里出现的 bug 先写回归测试，再修复；记录根因而非只写现象。
- 测事务回滚、对象删除失败、重复任务、并发替换。

**验收**

- 数据库提交失败时旧视频仍有效；清理失败不会导致用户更新接口失败或数据丢失。
- `go test -race ./...` 通过。
- Git：`day20: make lesson replacement recoverable`。

### Day 21：课程基本信息、更新、列表与详情

**听课（约 62 分钟）**

- 课程管理 1～3：《课程基本信息》《课程更新和列表》《课程详情和列表微调》

**亲手敲**

- `internal/course/`：创建、更新、发布/下架、分页列表、详情。
- 状态转换使用明确规则，例如 draft -> published -> archived；非法转换返回领域错误。
- 列表支持状态/讲师/关键词过滤，详情包含目录与课时摘要。
- 测状态机、RBAC、重复标题策略、分页和不存在资源。

**验收**

- 普通教师只能改自己负责的课程，管理员按权限操作全部课程。
- Git：`day21: implement course lifecycle`。

### Day 22：课程目录管理一、二

**听课（约 53 分钟）**

- 课程管理 4～5：《课程目录管理一》《课程目录管理二》

**亲手敲**

- 将 Day 17 catalog 与 course 垂直链路完整接通 handler/service/repository。
- 实现批量新建与批量排序，事务中校验所有节点属于同一课程。
- 返回稳定的树形 DTO；解决空 children 输出 `null` 还是 `[]` 的 API 一致性。
- 用 3 层、20 个节点的 fixture 测树构造和排序。

**验收**

- 任一节点校验失败时批量操作全部回滚。
- Git：`day22: integrate transactional course catalogs`。

### Day 23：课程目录管理三、四

**听课（约 45 分钟）**

- 课程管理 6～7：《课程目录管理三》《课程目录管理四》

**亲手敲**

- 实现移动节点、删除策略（禁止删除非空目录或显式级联二选一）、并发更新保护。
- 使用版本号/更新时间做乐观锁，冲突返回 409。
- 写属性式/表驱动测试：移动后无环、节点不丢、顺序唯一、跨课程失败。
- 对目录热点查询建立/验证索引。

**验收**

- 两个并发更新不能静默覆盖；冲突可被客户端识别并重试。
- Git：`day23: protect catalog moves from conflicts`。

### Day 24：整体部署、运行与演示收尾

**听课（约 17 分钟）**

- 课程管理 8～9：《整体部署运行》《第三阶段整体演示收尾》
- 至此 8 套指定课程的所有分集全部看完。

**亲手敲**

- `scripts/run.ps1`、`.env.example`、数据库迁移和 `docs/deployment.md`；构建 Windows 二进制：`go build ./cmd/api`。
- 运行 `gofmt`、`go vet`、全部测试、race test；为核心包生成 coverage，低于 70% 的核心 service 优先补失败路径。
- `docs/demo.md` 写可重复演示脚本：登录 -> 权限 -> 课程 -> 目录 -> 上传 -> 课时 -> 查询 -> 替换 -> 清理。
- 从空数据库严格按文档部署一次，不使用 IDE 隐藏配置。
- Docker 只作为加分项：若以后安装，再补 Dockerfile/Compose，不占用本日核心验收。

**验收**

```powershell
go build ./cmd/api
go test ./...
go test -race ./...
go vet ./...
```

- 新数据库按文档可启动，演示链路完整跑通。
- Git：`day24: make coursehub deployable and demonstrable`。

## 第五阶段：闭卷重建与总验收（Day 25）

### Day 25：不看课程，从零完成一条垂直链路

**今天不看新视频。**

**上午/第一阶段：90 分钟闭卷挑战**

- 在 `projects/coursehub-rebuild/` 中不看旧实现，只参考自己写的需求，继续使用仓库根 module，重建一个最小服务：配置、路由、鉴权中间件、Course repository/service/handler、创建与查询接口、测试。
- 90 分钟后再对照主项目，列出遗漏和错误。不要为了“看起来完成”复制旧代码。

**下午/第二阶段：系统验收**

- 运行所有 build/test/race/vet；手工演示完整业务链。
- 主动制造 5 个故障：MySQL 断开、token 过期、无权限、对象存储失败、并发目录冲突，确认响应和日志能定位。
- 检查密钥泄漏、日志敏感数据、未处理错误、goroutine 泄漏、事务边界、越权。
- 完成 `docs/final-retrospective.md`：掌握项、薄弱项、三个最难 bug、下一阶段 30 天计划。

**最终验收问题（必须能脱离代码回答）**

1. goroutine 怎样退出，谁负责关闭 channel，如何发现数据竞争？
2. middleware 的执行顺序是什么，401 与 403 有什么区别？
3. handler/service/repository 各自负责什么，事务应放在哪一层？
4. JWT 如何撤销，验证码如何防重放，OAuth state 防什么？
5. RBAC 怎样避免水平和垂直越权？
6. 对象上传成功但数据库失败怎么办？替换视频时怎样避免误删？
7. 树形目录怎样防环、保证顺序并处理并发覆盖？
8. 线上出现 500，如何通过 request ID、日志和测试复现定位？

**最终 Git 提交**：`day25: complete closed-book rebuild and final audit`。

## 掉队处理规则

- 当天少于 2 小时：课程必须听完，只实现标为核心的主路径和一个失败测试，次日先补测试再开始新内容。
- 卡住超过 30 分钟：写下“期望/实际/最小复现/已尝试”，缩小到最小测试，不继续盲改。
- 最多允许累计 1 天欠账；绝不靠最后三天倍速刷视频补课。
- 外部短信、飞书、腾讯云没有账号或付费条件时，使用接口 + fake adapter + 合同测试完成；这比把业务代码绑死在 SDK 上更能证明真正掌握。

## 每周自测节点

- Day 4：能独立写并发安全的 HTTP 小服务。
- Day 10：能从空目录搭出分层服务、数据库、日志和鉴权。
- Day 16：三种登录和 RBAC 权限矩阵可自动测试。
- Day 20：目录、对象存储、课时替换和失败补偿可运行。
- Day 24：完整项目可从空数据库部署和演示。
- Day 25：闭卷重建暴露真实掌握程度，并形成下一步计划。
