package exercises

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type Middleware func(http.Handler) http.Handler

type requestIDKey struct{}

type LogEntry struct {
	RequestID string
	Method    string
	Path      string
	Status    int
	Latency   time.Duration
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

// Chain 的参数按“请求进入时的执行顺序”书写。
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

var requestSequence atomic.Uint64

func RequestIDMiddleware(generator func() string) Middleware {
	if generator == nil {
		generator = func() string {
			return fmt.Sprintf("req-%d", requestSequence.Add(1))
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
			if requestID == "" {
				requestID = generator()
			}
			w.Header().Set("X-Request-ID", requestID)
			ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequestID(r *http.Request) string {
	requestID, _ := r.Context().Value(requestIDKey{}).(string)
	return requestID
}

func AccessLogMiddleware(logf func(LogEntry)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			writer := &statusWriter{ResponseWriter: w}
			next.ServeHTTP(writer, r)
			if writer.status == 0 {
				writer.status = http.StatusOK
			}
			if logf != nil {
				logf(LogEntry{
					RequestID: RequestID(r),
					Method:    r.Method,
					Path:      r.URL.Path,
					Status:    writer.status,
					Latency:   time.Since(start),
				})
			}
		})
	}
}

func RecoverMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recover() != nil {
					writeAPIError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func BearerAuthMiddleware(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Fields(r.Header.Get("Authorization"))
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] != token {
				writeAPIError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "valid bearer token required")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
