package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func TestLoginMultibytePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	old := DB
	DB = db
	t.Cleanup(func() { DB = old; db.Close() })
	password := "中文密"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery(`(?i)select id, username, password, email, role, status, created_at from users where username = \?`).WithArgs("tester").WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "email", "role", "status", "created_at"}).AddRow(1, "tester", hash, "test@example.com", "user", "enabled", time.Now()))
	router := gin.New()
	router.POST("/login", login)
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"tester","password":"中文密"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), hash) || strings.Contains(response.Body.String(), password) {
		t.Fatal("password leaked in response")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
