package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/websocket"
)

func runWebSocketExample() {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		_ = conn.WriteMessage(messageType, []byte("服务器收到: "+string(message)))
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		fmt.Println("WebSocket 连接失败:", err)
		return
	}
	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
		fmt.Println("发送失败:", err)
		return
	}
	_, response, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	fmt.Println("WebSocket 响应:", string(response))
}
