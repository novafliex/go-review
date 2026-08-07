package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type RESTCourse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func newRESTHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /courses/{id}", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(RESTCourse{ID: 1, Title: "课程 " + r.PathValue("id")})
	})
	mux.HandleFunc("POST /courses", func(w http.ResponseWriter, r *http.Request) {
		var course RESTCourse
		if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		course.ID = 1
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(course)
	})
	return mux
}

func runRESTExample() {
	body := bytes.NewBufferString(`{"title":"Go Web"}`)
	request := httptest.NewRequest(http.MethodPost, "/courses", body)
	response := httptest.NewRecorder()
	newRESTHandler().ServeHTTP(response, request)
	fmt.Printf("创建状态: %d，响应: %s", response.Code, response.Body.String())
}
