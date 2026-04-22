package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestExtractToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	token := extractToken(req)
	if token != "test-token" {
		t.Errorf("extractToken = %s, want test-token", token)
	}
}

func TestExtractTokenEmpty(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	token := extractToken(req)
	if token != "" {
		t.Errorf("extractToken should return empty string, got %s", token)
	}
}

func TestNewAuthMiddleware(t *testing.T) {
	am := NewAuthMiddleware("secret", "HS256")
	if am.secret != "secret" {
		t.Errorf("secret = %s, want secret", am.secret)
	}
}

func TestValidateTokenValid(t *testing.T) {
	am := NewAuthMiddleware("secret", "HS256")

	claims := &Claims{
		UserID: "user123",
		Role:   "ADMIN",
		Plan:   "ENTERPRISE",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))

	validated, err := am.validateToken(tokenString)
	if err != nil {
		t.Errorf("validateToken failed: %v", err)
	}

	if validated.UserID != "user123" {
		t.Errorf("UserID = %s, want user123", validated.UserID)
	}
}

func TestValidateTokenExpired(t *testing.T) {
	am := NewAuthMiddleware("secret", "HS256")

	claims := &Claims{
		UserID: "user123",
		Role:   "ADMIN",
		Plan:   "ENTERPRISE",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))

	_, err := am.validateToken(tokenString)
	if err == nil {
		t.Error("validateToken should fail for expired token")
	}
}

func TestNewRBACMiddleware(t *testing.T) {
	rm := NewRBACMiddleware("ADMIN")
	if rm.requiredRole != "ADMIN" {
		t.Errorf("requiredRole = %s, want ADMIN", rm.requiredRole)
	}
}
