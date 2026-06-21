package middlewares

import (
	"backend/internal/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("========================================")
		fmt.Println("🔐 AUTHENTICATE MIDDLEWARE CALLED")
		fmt.Printf("📋 Method: %s, Path: %s\n", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		fmt.Printf("📋 Authorization Header: '%s'\n", authHeader)

		if authHeader == "" {
			fmt.Println("❌ ERROR: Missing authorization header")
			utils.ErrorJson(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) {
			fmt.Printf("❌ ERROR: Header doesn't start with 'Bearer ', got: '%s'\n", authHeader)
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		fmt.Printf("🎫 Token extracted (length: %d): '%s...'\n", len(token), token[:min(len(token), 30)])

		userID, err := utils.VerifyAccessToken(token, secret)
		if err != nil {
			fmt.Printf("❌ ERROR: Token verification failed: %v\n", err)
			utils.ErrorJson(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		fmt.Printf("✅ Token verified! UserID: %d\n", *userID)
		fmt.Printf("🔧 Setting UserID in context: %d\n", *userID)

		ctx := context.WithValue(r.Context(), UserIDKey, *userID)

		fmt.Println("✅ Authentication successful, calling next handler")
		fmt.Println("========================================")

		next(w, r.WithContext(ctx))
	}
}
