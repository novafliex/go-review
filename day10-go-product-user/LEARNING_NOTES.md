# Learning Notes

## 2026-05-04 getIDFromRequest 重构

### 疑问

为什么不把 `writeJSON` 也封装进 `getIDFromRequest`？

一开始想这样写：

```go
func getIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "id is required",
			Data:    nil,
		})
		return 0, false
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid id",
			Data:    nil,
		})
		return 0, false
	}

	return id, true
}
```

### 问题

这样虽然能运行，但不推荐。

因为 `getIDFromRequest` 的职责应该是：

```text
从请求中解析 id
```

而不是：

```text
解析 id + 返回 JSON 响应
```

如果把 `writeJSON` 放进去，返回响应的逻辑就被藏在工具函数里了，后面看业务代码时不够清楚。

### 推荐写法

#### 1. `getIDFromRequest` 只负责解析 id

```go
func getIDFromRequest(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		return 0, fmt.Errorf("id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}

	return id, nil
}
```

#### 2. `writeError` 负责返回错误 JSON

```go
func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, Response{
		Code:    statusCode,
		Message: message,
		Data:    nil,
	})
}
```

#### 3. 业务函数中组合使用

```go
id, err := getIDFromRequest(r)
if err != nil {
	writeError(w, http.StatusBadRequest, err.Error())
	return
}
```

### 结论

```text
getIDFromRequest：负责解析请求参数
writeError / writeJSON：负责返回 JSON 响应
updateProduct / deleteProduct：负责业务流程控制
```

封装不是把代码全部塞进一个函数，而是把**同一类职责**的代码放到一起。

## 2026-05-05 V2 到 V3：接入 MySQL 需要的导包

### 疑问

V2 到 V3 接入 MySQL 时，需要导入什么包？

```bash
go get github.com/go-sql-driver/mysql
go get github.com/joho/godotenv
```

### 问题

```go
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)
```

`github.com/go-sql-driver/mysql` 是 MySQL 驱动，让 Go 可以使用：

```go
sql.Open("mysql", dsn)
```

前面的 `_` 表示只注册驱动，不直接调用这个包里的函数。

`github.com/joho/godotenv` 用来读取 `.env` 文件，把数据库配置加载到环境变量里。

### 推荐写法

数据库账号和密码不要写死在代码里：

```go
dsn := "root:123456@tcp(127.0.0.1:3306)/product_manager"
```

更推荐写到 `.env`：

```text
DB_USER=root
DB_PASSWORD=你的MySQL密码
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=product_manager
```

然后在 Go 代码里用 `os.Getenv` 读取，再拼接 DSN。

### 结论

```text
mysql 驱动包：让 Go 可以连接 MySQL
godotenv：读取 .env 配置
.env：保存数据库账号、密码等敏感信息
```

数据库配置不要写死在代码里，放到 `.env` 更安全，也方便以后修改。
