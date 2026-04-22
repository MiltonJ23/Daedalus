package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/MILTONJ23/Daedalus/Shared/errors"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Plan   string `json:"plan"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	secret    string
	algorithm string
}

func NewAuthMiddleware(secret string, algorithm string) *AuthMiddleware {
	return &AuthMiddleware{
		secret:    secret,
		algorithm: algorithm,
	}
}

func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			writeError(w, errors.NewUnauthorized("missing token"))
			return
		}

		claims, err := am.validateToken(tokenString)
		if err != nil {
			writeError(w, errors.NewUnauthorized(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)
		ctx = context.WithValue(ctx, "plan", claims.Plan)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != am.algorithm {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

type RBACMiddleware struct {
	requiredRole string
}

func NewRBACMiddleware(role string) *RBACMiddleware {
	return &RBACMiddleware{requiredRole: role}
}

func (rm *RBACMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != rm.requiredRole {
			writeError(w, errors.NewForbidden("insufficient permissions"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

type LoggingMiddleware struct {
	handler http.Handler
}

func NewLoggingMiddleware(handler http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{handler: handler}
}

func (lm *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lm.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	_ = fmt.Sprintf("%s %s - %v", r.Method, r.RequestURI, duration)
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func writeError(w http.ResponseWriter, err *errors.DaedalusError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus())
}
