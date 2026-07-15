# Go Web Product Manager

一个基于 **Go + Gin + MySQL + Vue 3 CDN** 的商品管理 CRUD 项目。

项目从 Go Web 基础练习逐步演进到完整的前后端联调版本：后端提供 RESTful API，数据持久化到 MySQL；前端使用单文件 `frontend/index.html`，通过 Vue 3 CDN 直接调用 Gin 接口，实现商品的查询、新增、编辑、删除、搜索筛选、深浅色模式、Toast 提示和删除确认弹窗。

## 项目功能

- 商品列表查询
- 根据商品 ID 查询单个商品
- 新增商品
- 修改商品
- 删除商品
- MySQL 数据持久化
- 统一 JSON 响应格式
- 基础参数校验和错误响应
- 简单 CORS 中间件，支持本地前端访问后端接口
- Vue 3 CDN 单文件前端
- 现代 SaaS Dashboard 风格页面
- 深色 / 浅色模式切换
- 商品搜索和最高价格筛选
- 商品列表内部滚动和 sticky 表头
- 删除确认弹窗
- Toast 操作反馈
- 表格展示连续“序号”，同时保留数据库真实“商品 ID”

## 技术栈

### 后端

- Go
- Gin
- MySQL
- `database/sql`
- `go-sql-driver/mysql`
- `godotenv`

### 前端

- HTML
- CSS
- Vue 3 CDN
- Fetch API

本项目没有引入 Vue CLI、Vite 或任何 UI 框架，前端保持单文件 HTML。

## 项目结构

```text
go-web-product-manager/
├── main.go              # 注册路由、CORS 中间件并启动 Gin 服务
├── database.go          # 读取环境变量并初始化 MySQL 连接
├── handler.go           # 商品 CRUD 接口处理函数
├── model.go             # Product 和 Response 结构体
├── frontend/
│   └── index.html       # Vue 3 CDN 单文件前端页面
├── LEARNING_NOTES.md    # 学习笔记
├── NOTE_TEMPLATE.md     # 笔记模板
├── go.mod
├── go.sum
└── README.md
```

## 数据库准备

创建数据库和商品表：

```sql
CREATE DATABASE IF NOT EXISTS product_manager
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE product_manager;

CREATE TABLE IF NOT EXISTS products (
  id INT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  price DECIMAL(10, 2) NOT NULL
);
```

可以插入一些初始数据：

```sql
INSERT INTO products (id, name, price) VALUES
  (1, '手机', 5200.00),
  (2, '电脑', 8888.88),
  (3, '相机', 3599.00);
```

说明：`id` 是数据库主键，删除商品后不会自动重排。前端表格中的“序号”只是当前筛选结果的展示序号，“商品 ID”才是数据库真实 ID。

## 环境变量

项目通过 `.env` 或系统环境变量读取 MySQL 连接字符串：

```env
MYSQL_DSN=用户名:密码@tcp(127.0.0.1:3306)/product_manager?charset=utf8mb4&parseTime=True&loc=Local
```

`.env` 用于本地开发，不要提交真实数据库账号和密码。

## 启动后端

安装依赖：

```bash
go mod tidy
```

启动 Gin 服务：

```bash
go run .
```

启动成功后，后端服务地址为：

```text
http://localhost:8080
```

## 启动前端

直接用浏览器打开：

```text
frontend/index.html
```

前端会请求：

```text
http://localhost:8080/products
```

后端已经在 `main.go` 中添加了简单 CORS 中间件，方便本地 HTML 页面直接访问 Gin 接口。

## 接口列表

| 方法     | 路径            | 说明                 |
| -------- | --------------- | -------------------- |
| `GET`    | `/products`     | 获取商品列表         |
| `GET`    | `/products/:id` | 根据 ID 获取单个商品 |
| `POST`   | `/products`     | 新增商品             |
| `PUT`    | `/products/:id` | 修改商品             |
| `DELETE` | `/products/:id` | 删除商品             |

## 统一响应格式

接口统一返回：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

示例：商品列表响应

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "手机",
      "price": 5200
    }
  ]
}
```

## 请求示例

### 获取商品列表

```bash
curl http://localhost:8080/products
```

### 获取单个商品

```bash
curl http://localhost:8080/products/1
```

### 新增商品

新增商品时需要传入 `id`、`name`、`price`：

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d "{\"id\":5,\"name\":\"键盘\",\"price\":299}"
```

### 修改商品

编辑商品时 ID 来自路径参数，不通过表单修改：

```bash
curl -X PUT http://localhost:8080/products/5 \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"机械键盘\",\"price\":699}"
```

### 删除商品

```bash
curl -X DELETE http://localhost:8080/products/5
```

## 常见错误

| 状态码 | message                             | 说明                                        |
| ------ | ----------------------------------- | ------------------------------------------- |
| `400`  | `invalid id`                        | 路径中的 ID 不是数字                        |
| `400`  | `id is required`                    | 缺少路径 ID                                 |
| `400`  | `invalid product`                   | 请求体不是合法 JSON，或不能绑定到 `Product` |
| `400`  | `product id/name/price is required` | 新增商品时缺少必要字段                      |
| `400`  | `product name/price is required`    | 修改商品时缺少必要字段                      |
| `404`  | `not found`                         | 商品不存在                                  |
| `409`  | `product already exists or failed`  | 新增商品 ID 已存在，或插入失败              |
| `500`  | `query products failed`             | 查询商品列表失败                            |
| `500`  | `update product failed`             | 修改商品失败                                |
| `500`  | `delete product failed`             | 删除商品失败                                |

## 前端页面说明

`frontend/index.html` 是一个完整的 Vue 3 CDN 单文件页面，包含：

- 页面 HTML
- 页面 CSS
- Vue 业务逻辑
- Fetch 接口调用

前端接口调用对应关系：

| 前端函数        | 后端接口                  |
| --------------- | ------------------------- |
| `loadProducts`  | `GET /products`           |
| `createProduct` | `POST /products`          |
| `updateProduct` | `PUT /products/:id`       |
| `deleteProduct` | `DELETE /products/:id`    |

前端表格中：

- “序号”根据当前 `filteredProducts` 从 1 开始生成。
- “商品 ID”显示数据库真实 `product.id`。
- 编辑和删除始终使用 `product.id`，不会使用序号。
- 搜索筛选后，序号会根据当前显示结果重新连续排列。

## 学习收获

这个项目覆盖了一个小型 Go Web CRUD 应用的完整路径：

- 从 HTTP 路由到 RESTful API
- 从内存数据到 MySQL 持久化
- 从后端接口到前端页面联调
- 从普通表格到可用的 Dashboard 管理界面
- 从零散响应到统一 JSON 响应结构
- 从基础功能到错误处理、CORS、表单校验和用户反馈

## 当前状态

项目已经完成：

- 后端 CRUD 接口完成
- MySQL 持久化完成
- 前端真实接口接入完成
- Dashboard 页面优化完成
- 深浅色模式完成
- 搜索筛选完成
- 删除确认和 Toast 提示完成
- 本地 CORS 支持完成

这是一个适合作为 Go Web 入门阶段总结的小型全栈 CRUD 项目。
