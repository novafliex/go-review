package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// password 存的是加密后的密码
type User struct {
	Username string
	Password string
}

func Register(username string, password string) User {
	// 明文密码转为加密密码
	// bcrypt.GenerateFromPassword() 第一个参数是明文密码(bcrypt要求类型是[]byte) 第二个是加密强度
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 为什么是panic(err)
	if err != nil {
		panic(err)
	}
	// hashedPassword 是[]byte类型
	user := User{
		Username: username,
		Password: string(hashedPassword),
	}
	return user
}

func Login(user User, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), // 数据库里加密的密码
		[]byte(password),      // 用户输入的密码
	)
	if err != nil {
		return false
	}
	return true
}

func main() {
	user := Register(
		"admin",
		"123456",
	)

	fmt.Println("username:", user.Username)
	fmt.Println("database password:", user.Password)

	// 模拟登录正确密码
	result := Login(
		user,
		"123456",
	)

	fmt.Println("第一次登录:", result)

	// 模拟登录错误密码
	result = Login(
		user,
		"111111",
	)

	fmt.Println("第二次登录:", result)

}
