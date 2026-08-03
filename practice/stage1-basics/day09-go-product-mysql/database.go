package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// 一个全局变量 用来保存数据库连接池
var DB *sql.DB

func InitDB() error {
	// 读取项目里的 .env 文件
	// godotenv.Load() 无论成功还是失败，都会返回一个 error 类型的值：加载成功：返回 nil ,加载失败：返回具体错误，赋值给 err
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found, use system environment variables, the err is :", err)
	}

	// 获取 mysql 连接信息  只想拿值：os.Getenv 还想知道存不存在：os.LookupEnv
	// os.Getenv("MYSQL_DSN")从当前程序的环境变量中，读取名字叫 MYSQL_DSN 的值
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return fmt.Errorf("MYSQL_DSN is empty")
	}

	// 创建 MySQL 创建连接池对象，不一定马上真正连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 真正测试 MySQL 能不能连接
	err = db.Ping()
	if err != nil {
		db.Close()
		return err
	}

	// 把连接成功的局部变量 db，保存到全局变量 DB，方便其他接口使用。
	DB = db
	fmt.Println("MYSQL connected successfully")
	return nil
}

//1. 声明全局变量 DB，用来保存数据库连接池
// 2. godotenv.Load() 加载 .env 文件
// 3. os.Getenv("MYSQL_DSN") 读取 MySQL 连接字符串
// 4. sql.Open("mysql", dsn) 创建连接池对象，先保存到局部变量 db
// 5. db.Ping() 真正测试数据库能否连接
// 6. 如果 Ping 失败，db.Close() 并返回错误
// 7. DB = db，把同一个连接池保存到全局变量
// 8. return nil，表示初始化成功
