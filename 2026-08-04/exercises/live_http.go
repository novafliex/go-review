package exercises

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
)

var (
	ErrInvalidCourse = errors.New("course title is required")
	ErrCourseExists  = errors.New("course title already exists")
)

type LiveCourse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type LiveEvent struct {
	Type   string     `json:"type"`
	Course LiveCourse `json:"course"`
}

type Publisher interface {
	Publish(LiveEvent)
}

type LiveCourseStore struct {
	mu        sync.RWMutex
	nextID    int
	courses   []LiveCourse
	publisher Publisher
}

func NewLiveCourseStore(publisher Publisher) *LiveCourseStore {
	return &LiveCourseStore{nextID: 1, publisher: publisher}
}

func (s *LiveCourseStore) List() []LiveCourse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]LiveCourse(nil), s.courses...)
}

func (s *LiveCourseStore) Create(title string) (LiveCourse, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return LiveCourse{}, ErrInvalidCourse
	}

	s.mu.Lock()
	for _, course := range s.courses {
		if strings.EqualFold(course.Title, title) {
			s.mu.Unlock()
			return LiveCourse{}, ErrCourseExists
		}
	}
	course := LiveCourse{ID: s.nextID, Title: title}
	s.nextID++
	s.courses = append(s.courses, course)
	s.mu.Unlock()

	if s.publisher != nil {
		s.publisher.Publish(LiveEvent{Type: "course.created", Course: course})
	}
	return course, nil
}

func NewLiveCourseHTTPHandler(store *LiveCourseStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeLiveJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /courses", func(w http.ResponseWriter, _ *http.Request) {
		writeLiveJSON(w, http.StatusOK, store.List())
	})
	mux.HandleFunc("POST /courses", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title string `json:"title"`
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&input); err != nil {
			writeLiveJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
			return
		}

		course, err := store.Create(input.Title)
		switch {
		case errors.Is(err, ErrInvalidCourse):
			writeLiveJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, ErrCourseExists):
			writeLiveJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		case err != nil:
			writeLiveJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		default:
			writeLiveJSON(w, http.StatusCreated, course)
		}
	})
	return mux
}

func writeLiveJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
