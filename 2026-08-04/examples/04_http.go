package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type HTTPStudent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func runHTTPExample() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /students", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]HTTPStudent{{ID: 1, Name: "张三"}})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := http.Get(server.URL + "/students")
	if err != nil {
		fmt.Println("HTTP 请求失败:", err)
		return
	}
	defer response.Body.Close()

	var students []HTTPStudent
	if err := json.NewDecoder(response.Body).Decode(&students); err != nil {
		fmt.Println("解析响应失败:", err)
		return
	}
	fmt.Printf("HTTP 状态: %d，学生: %+v\n", response.StatusCode, students)
}
