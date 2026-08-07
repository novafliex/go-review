package exercises

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type recordingPublisher struct {
	events []LiveEvent
}

func (p *recordingPublisher) Publish(event LiveEvent) {
	p.events = append(p.events, event)
}

func TestCourseHTTPHandler(t *testing.T) {
	publisher := &recordingPublisher{}
	handler := NewLiveCourseHTTPHandler(NewLiveCourseStore(publisher))

	t.Run("health", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", response.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBufferString("{"))
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", response.Code)
		}
	})

	t.Run("create and list", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBufferString(`{"title":"Go 并发"}`))
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201: %s", response.Code, response.Body.String())
		}
		if len(publisher.events) != 1 || publisher.events[0].Type != "course.created" {
			t.Fatalf("unexpected events: %+v", publisher.events)
		}

		request = httptest.NewRequest(http.MethodGet, "/courses", nil)
		response = httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		var courses []LiveCourse
		if err := json.NewDecoder(response.Body).Decode(&courses); err != nil {
			t.Fatal(err)
		}
		if len(courses) != 1 || courses[0].Title != "Go 并发" {
			t.Fatalf("unexpected courses: %+v", courses)
		}
	})
}
