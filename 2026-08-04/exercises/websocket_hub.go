package exercises

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	conn *websocket.Conn
	send chan LiveEvent
}

type Hub struct {
	mu       sync.RWMutex
	clients  map[*wsClient]struct{}
	upgrader websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*wsClient]struct{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: sameOrigin,
		},
	}
}

func sameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Host == r.Host
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &wsClient{conn: conn, send: make(chan LiveEvent, 8)}

	h.mu.Lock()
	h.clients[client] = struct{}{}
	h.mu.Unlock()

	go h.writePump(client)
	defer h.remove(client)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) Publish(event LiveEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.send <- event:
		default:
			// 慢客户端不能阻塞所有广播；写满后让连接退出。
			_ = client.conn.Close()
		}
	}
}

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *Hub) writePump(client *wsClient) {
	for event := range client.send {
		if err := client.conn.WriteJSON(event); err != nil {
			_ = client.conn.Close()
			return
		}
	}
}

func (h *Hub) remove(client *wsClient) {
	h.mu.Lock()
	if _, exists := h.clients[client]; exists {
		delete(h.clients, client)
		close(client.send)
	}
	h.mu.Unlock()
	_ = client.conn.Close()
}
