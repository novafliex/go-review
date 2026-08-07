package exercises

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition was not met")
}

func TestHubBroadcastAndDisconnect(t *testing.T) {
	hub := NewHub()
	server := httptest.NewServer(http.HandlerFunc(hub.ServeWS))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { return hub.Count() == 1 })

	want := LiveEvent{Type: "course.created", Course: LiveCourse{ID: 1, Title: "Go 网络"}}
	hub.Publish(want)
	var got LiveEvent
	if err := conn.ReadJSON(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	_ = conn.Close()
	waitFor(t, func() bool { return hub.Count() == 0 })
}
