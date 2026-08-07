package exercises

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	errNotFound = errors.New("course not found")
	errConflict = errors.New("course title already exists")
	errInvalid  = errors.New("course title is required")
)

type Course struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type APIError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type successResponse struct {
	Data any `json:"data"`
}

type courseRepository struct {
	mu      sync.RWMutex
	nextID  int
	courses map[int]Course
}

func newCourseRepository() *courseRepository {
	return &courseRepository{nextID: 1, courses: make(map[int]Course)}
}

func (r *courseRepository) list() []Course {
	r.mu.RLock()
	defer r.mu.RUnlock()
	courses := make([]Course, 0, len(r.courses))
	for _, course := range r.courses {
		courses = append(courses, course)
	}
	sort.Slice(courses, func(i, j int) bool { return courses[i].ID < courses[j].ID })
	return courses
}

func (r *courseRepository) get(id int) (Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	course, exists := r.courses[id]
	if !exists {
		return Course{}, errNotFound
	}
	return course, nil
}

func (r *courseRepository) create(title string) (Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	title = strings.TrimSpace(title)
	if title == "" {
		return Course{}, errInvalid
	}
	if r.titleExists(title, 0) {
		return Course{}, errConflict
	}
	course := Course{ID: r.nextID, Title: title}
	r.nextID++
	r.courses[course.ID] = course
	return course, nil
}

func (r *courseRepository) update(id int, title string) (Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.courses[id]; !exists {
		return Course{}, errNotFound
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return Course{}, errInvalid
	}
	if r.titleExists(title, id) {
		return Course{}, errConflict
	}
	course := Course{ID: id, Title: title}
	r.courses[id] = course
	return course, nil
}

func (r *courseRepository) delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.courses[id]; !exists {
		return errNotFound
	}
	delete(r.courses, id)
	return nil
}

func (r *courseRepository) titleExists(title string, exceptID int) bool {
	for id, course := range r.courses {
		if id != exceptID && strings.EqualFold(course.Title, title) {
			return true
		}
	}
	return false
}

type courseAPI struct {
	repository *courseRepository
}

func NewAPI(token string, logf func(LogEntry)) http.Handler {
	api := &courseAPI{repository: newCourseRepository()}
	auth := BearerAuthMiddleware(token)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /courses", api.listCourses)
	mux.Handle("POST /courses", auth(http.HandlerFunc(api.createCourse)))
	mux.HandleFunc("GET /courses/{id}", api.getCourse)
	mux.Handle("PUT /courses/{id}", auth(http.HandlerFunc(api.updateCourse)))
	mux.Handle("DELETE /courses/{id}", auth(http.HandlerFunc(api.deleteCourse)))

	// 无方法模式作为兜底，让 404/405 也保持统一 JSON 格式。
	mux.HandleFunc("/courses", api.methodNotAllowed)
	mux.HandleFunc("/courses/{id}", api.methodNotAllowed)
	mux.HandleFunc("/", api.notFound)

	return Chain(
		mux,
		RequestIDMiddleware(nil),
		AccessLogMiddleware(logf),
		RecoverMiddleware(),
	)
}

func (api *courseAPI) listCourses(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, successResponse{Data: api.repository.list()})
}

func (api *courseAPI) getCourse(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	course, err := api.repository.get(id)
	if err != nil {
		writeRepositoryError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, successResponse{Data: course})
}

func (api *courseAPI) createCourse(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	course, err := api.repository.create(input.Title)
	if err != nil {
		writeRepositoryError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, successResponse{Data: course})
}

func (api *courseAPI) updateCourse(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input struct {
		Title string `json:"title"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	course, err := api.repository.update(id, input.Title)
	if err != nil {
		writeRepositoryError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, successResponse{Data: course})
}

func (api *courseAPI) deleteCourse(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := api.repository.delete(id); err != nil {
		writeRepositoryError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (api *courseAPI) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeAPIError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
}

func (api *courseAPI) notFound(w http.ResponseWriter, r *http.Request) {
	writeAPIError(w, r, http.StatusNotFound, "NOT_FOUND", "route not found")
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeAPIError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid course id")
		return 0, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeAPIError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON body")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeAPIError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "body must contain one JSON object")
		return false
	}
	return true
}

func writeRepositoryError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errInvalid):
		writeAPIError(w, r, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
	case errors.Is(err, errNotFound):
		writeAPIError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, errConflict):
		writeAPIError(w, r, http.StatusConflict, "CONFLICT", err.Error())
	default:
		writeAPIError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func writeAPIError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	writeJSON(w, status, APIError{Code: code, Message: message, RequestID: RequestID(r)})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
