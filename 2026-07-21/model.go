//结构体定义

package main

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	ID    int     `json:"id"`
}
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role,omitempty"`
	CreatedAt string `json:"created_at"`
}
