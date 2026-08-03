package main

import "golang.org/x/crypto/bcrypt"

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
