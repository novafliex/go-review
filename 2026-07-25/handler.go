package main

import (
	"database/sql"
	"errors"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	// tips: 不可以清空密码的空格 而是设定密码中不得包含任何空白字符

	// check username
	if req.Username == "" {
		ErrorResponse(c, 400, "username is required")
		return
	}
	if len([]rune(req.Username)) < 3 || len([]rune(req.Username)) > 50 {
		ErrorResponse(c, 400, "username must be between 3 and 50 characters")
		return
	}

	// check email
	if req.Email == "" {
		ErrorResponse(c, 400, "email is required")
		return
	}
	if len([]rune(req.Email)) > 100 {
		ErrorResponse(c, 400, "email cannot exceed 100 characters")
		return
	}

	// check password
	if req.Password == "" {
		ErrorResponse(c, 400, "password is required")
		return
	}
	// it is prohibited for passwords to contain common white space characters such as spaces, tabs, and newlines
	if strings.IndexFunc(req.Password, unicode.IsSpace) != -1 {
		ErrorResponse(c, 400, "password cannot contain spaces")
		return
	}

	if len([]byte(req.Password)) < 6 {
		ErrorResponse(c, 400, "password must be at least 6 bytes")
		return
	}

	// bcrypt 最多处理 72 字节
	if len([]byte(req.Password)) > 72 {
		ErrorResponse(c, 400, "password cannot exceed 72 bytes")
		return
	}

	// 6. 检查用户名是否重复
	var existingID int

	err := DB.QueryRow(
		"SELECT id FROM users WHERE username = ?",
		req.Username,
	).Scan(&existingID)

	if err == nil {
		ErrorResponse(c, 409, "username already exists")
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 500, "check username failed")
		return
	}

	// 7. 检查邮箱是否重复
	err = DB.QueryRow(
		"SELECT id FROM users WHERE email = ?",
		req.Email,
	).Scan(&existingID)

	if err == nil {
		ErrorResponse(c, 409, "email already exists")
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 500, "check email failed")
		return
	}
	// 8. 加密密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		ErrorResponse(c, 500, "hash password failed")
		return
	}

	// 9. 设置默认角色和状态
	role := "user"
	status := "enabled"

	// 10. 插入数据库
	result, err := DB.Exec(
		`INSERT INTO users
				(username, password, email, role, status)
			 VALUES (?, ?, ?, ?, ?)`,
		req.Username,
		hashedPassword,
		req.Email,
		role,
		status,
	)
	if err != nil {
		ErrorResponse(c, 500, "register user failed")
		return
	}

	// 11. 获取新增用户 ID
	id, err := result.LastInsertId()
	if err != nil {
		ErrorResponse(c, 500, "get last insert id failed")
		return
	}

	// 12. 查询并返回新增用户，不查询密码
	var user User

	err = DB.QueryRow(
		`SELECT id, username, email, role, status, created_at
			 FROM users
			 WHERE id = ?`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
	)
	if err != nil {
		ErrorResponse(c, 500, "query registered user failed")
		return
	}

	JSONResponse(c, 201, "user registered successfully", user)
}

func login(c *gin.Context) {
	// 接受登录信息
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, 400, "invaild request body")
		return
	}

	// 清理用户名
	req.Username = strings.TrimSpace(req.Username)

	// 参数校验
	if req.Username == "" {
		ErrorResponse(c, 400, "username is required")
		return
	}
	if req.Password == "" {
		ErrorResponse(c, 400, "password is required")
		return
	}
	if len([]rune(req.Username)) < 3 || len([]rune(req.Username)) > 50 {
		ErrorResponse(c, 401, "username or password is incorrect")
		return
	}
	if strings.IndexFunc(req.Password, unicode.IsSpace) != -1 {
		ErrorResponse(c, 401, "username or password is incorrect")
		return
	}
	passwordLength := len([]byte(req.Password))
	if passwordLength < 6 || passwordLength > 72 {
		ErrorResponse(c, 401, "username or password is incorrect")
		return
	}

	// 根据用户名查用户
	var user User

	err := DB.QueryRow("select id, username, password, email, role, status, created_at from users where username = ?", req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
	)
	// the username does not exist
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(c, 401, "username or password is incorrect")
		return
	}
	// other database errors
	if err != nil {
		ErrorResponse(c, 500, "query user failed")
		return
	}

	// 比较密码
	if !CheckPassword(user.Password, req.Password) {
		ErrorResponse(c, 401, "username or password is incorrect")
		return
	}

	// 检查账号状态
	if user.Status != "enabled" {
		ErrorResponse(c, 403, "user account is disabled")
		return
	}

	// 清空密码字段
	user.Password = ""

	// 返回登录成功
	JSONResponse(c, 200, "login successful", user)
}
