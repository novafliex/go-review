package main

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	contextUserIDKey = "userID"
	contextRoleKey   = "role"
)

// ParseToken verifies a JWT and returns its claims.
func ParseToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid || claims.UserID <= 0 || claims.Role == "" {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// AuthMiddleware protects routes that require a valid Bearer token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			ErrorResponse(c, http.StatusUnauthorized, "authorization header must be Bearer token")
			c.Abort()
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(contextUserIDKey, claims.UserID)
		c.Set(contextRoleKey, claims.Role)
		c.Next()
	}
}
