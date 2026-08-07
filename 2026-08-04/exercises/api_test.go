package exercises

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

const testToken = "coursehub-token"

func request(handler http.Handler, method, path, body string, authenticated bool) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if authenticated {
		req.Header.Set("Authorization", "Bearer "+testToken)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func TestCourseCRUD(t *testing.T) {
	handler := NewAPI(testToken, nil)

	response := request(handler, http.MethodPost, "/courses", `{"title":"Go Web"}`, false)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated create got %d, want 401", response.Code)
	}

	response = request(handler, http.MethodPost, "/courses", `{"title":"Go Web"}`, true)
	if response.Code != http.StatusCreated {
		t.Fatalf("create got %d, want 201: %s", response.Code, response.Body.String())
	}

	response = request(handler, http.MethodPost, "/courses", `{"title":"go web"}`, true)
	if response.Code != http.StatusConflict {
		t.Fatalf("duplicate got %d, want 409", response.Code)
	}

	response = request(handler, http.MethodGet, "/courses", "", false)
	if response.Code != http.StatusOK || !bytes.Contains(response.Body.Bytes(), []byte("Go Web")) {
		t.Fatalf("unexpected list response: %d %s", response.Code, response.Body.String())
	}

	response = request(handler, http.MethodGet, "/courses/1", "", false)
	if response.Code != http.StatusOK {
		t.Fatalf("get got %d, want 200", response.Code)
	}

	response = request(handler, http.MethodPut, "/courses/1", `{"title":"Go 网络编程"}`, true)
	if response.Code != http.StatusOK || !bytes.Contains(response.Body.Bytes(), []byte("Go 网络编程")) {
		t.Fatalf("unexpected update response: %d %s", response.Code, response.Body.String())
	}

	response = request(handler, http.MethodDelete, "/courses/1", "", true)
	if response.Code != http.StatusNoContent {
		t.Fatalf("delete got %d, want 204", response.Code)
	}

	response = request(handler, http.MethodGet, "/courses/1", "", false)
	if response.Code != http.StatusNotFound {
		t.Fatalf("deleted course got %d, want 404", response.Code)
	}
}

func TestUnifiedErrorsAndRequestID(t *testing.T) {
	var logs []LogEntry
	handler := NewAPI(testToken, func(entry LogEntry) { logs = append(logs, entry) })

	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBufferString("{"))
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("X-Request-ID", "fixed-id")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400", response.Code)
	}
	var apiError APIError
	if err := json.NewDecoder(response.Body).Decode(&apiError); err != nil {
		t.Fatal(err)
	}
	if apiError.Code != "INVALID_REQUEST" || apiError.RequestID != "fixed-id" {
		t.Fatalf("unexpected error: %+v", apiError)
	}
	if len(logs) != 1 || logs[0].Status != http.StatusBadRequest || logs[0].RequestID != "fixed-id" {
		t.Fatalf("unexpected logs: %+v", logs)
	}

	response = request(handler, http.MethodPatch, "/courses", `{}`, false)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("method got %d, want 405", response.Code)
	}
	response = request(handler, http.MethodGet, "/missing", "", false)
	if response.Code != http.StatusNotFound {
		t.Fatalf("route got %d, want 404", response.Code)
	}
}

func TestRecoverMiddleware(t *testing.T) {
	handler := Chain(
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }),
		RequestIDMiddleware(func() string { return "panic-id" }),
		RecoverMiddleware(),
	)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/panic", nil))
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("got %d, want 500", response.Code)
	}
	var apiError APIError
	_ = json.NewDecoder(response.Body).Decode(&apiError)
	if apiError.Code != "INTERNAL_ERROR" || apiError.RequestID != "panic-id" {
		t.Fatalf("unexpected error: %+v", apiError)
	}
}

func TestMiddlewareOrder(t *testing.T) {
	var order []string
	track := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+" before")
				next.ServeHTTP(w, r)
				order = append(order, name+" after")
			})
		}
	}
	final := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		order = append(order, "handler")
	})
	Chain(final, track("A"), track("B")).ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", nil),
	)
	want := []string{"A before", "B before", "handler", "B after", "A after"}
	if !slices.Equal(order, want) {
		t.Fatalf("got %v, want %v", order, want)
	}
}
