package main

import (
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type passwordHash struct{ plain string }

func (p passwordHash) Match(value driver.Value) bool {
	hash, ok := value.(string)
	return ok && bcrypt.CompareHashAndPassword([]byte(hash), []byte(p.plain)) == nil
}

func TestUserWrites(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, update := range []bool{false, true} {
		name := "create"
		if update {
			name = "update"
		}
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			old := DB
			DB = db
			t.Cleanup(func() { DB = old; db.Close() })
			router := gin.New()
			method, path, want := http.MethodPost, "/users", http.StatusCreated
			if update {
				method, path, want = http.MethodPut, "/users/1", http.StatusOK
				router.PUT("/users/:id", updateUser)
				mock.ExpectExec(`update users set username = \?, password = \?, email = \?, phone = \? where id = \?`).WithArgs("tester", passwordHash{"secret123"}, "test@example.com", "123", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			} else {
				router.POST("/users", createUser)
				mock.ExpectExec(`insert into users`).WithArgs("tester", passwordHash{"secret123"}, "test@example.com", "123", "user").WillReturnResult(sqlmock.NewResult(1, 1))
			}
			request := httptest.NewRequest(method, path, strings.NewReader(`{"username":"tester","password":"secret123","email":"test@example.com","phone":"123","role":"admin"}`))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != want {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			if strings.Contains(response.Body.String(), `"password"`) || strings.Contains(response.Body.String(), "admin") {
				t.Fatalf("unsafe response: %s", response.Body.String())
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
