package main

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	ID    int     `json:"id"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
