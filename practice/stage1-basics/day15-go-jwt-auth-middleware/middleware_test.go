package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateToken(7, "user")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{name: "valid token", header: "Bearer " + token, wantStatus: http.StatusOK},
		{name: "missing header", wantStatus: http.StatusUnauthorized},
		{name: "invalid token", header: "Bearer invalid", wantStatus: http.StatusUnauthorized},
		{name: "wrong scheme", header: "Basic " + token, wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
				userID, exists := c.Get(contextUserIDKey)
				if !exists || userID != 7 {
					t.Fatalf("context user ID = %v, exists = %v", userID, exists)
				}
				c.Status(http.StatusOK)
			})

			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.header != "" {
				request.Header.Set("Authorization", tt.header)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body = %s", response.Code, tt.wantStatus, response.Body.String())
			}
		})
	}
}

func TestParseTokenRejectsMissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if _, err := ParseToken("token"); err == nil {
		t.Fatal("ParseToken() expected an error")
	}
}
