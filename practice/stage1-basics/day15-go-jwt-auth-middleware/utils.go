package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// claims 表示jwt中保存的用户信息
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

// 注册时把明文密码转换为加密密码
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// 登陆时比较数据库密码和用户输入的密码
func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), // 数据库中的加密密码
		[]byte(password),       // 用户输入的明文密码
	)

	return err == nil
}

// GenerateToken 登录成功后生成jwt
func GenerateToken(userID int, role string) (string, error) {
	// 从环境变量中读取密钥
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is empty")
	}
	// 获取当前时间
	now := time.Now()

	//设置jwt中保存的数据
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// token 签发时间
			IssuedAt: jwt.NewNumericDate(now),
			// token 过期时间
			ExpiresAt: jwt.NewNumericDate(
				now.Add(24 * time.Hour),
			),
		},
	}

	// 4. 创建一个使用 HS256 算法的 Token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// 5. 使用 JWT_SECRET 对 Token 签名
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
