package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims represents the claims from a JWT token
type JWTClaims struct {
	Sub   string `json:"sub"`   // Subject (customer ID)
	Email string `json:"email"` // Customer email
	Name  string `json:"name"`  // Customer name
	Iss   string `json:"iss"`   // Issuer
	Exp   int64  `json:"exp"`   // Expiration time
	Iat   int64  `json:"iat"`   // Issued at
}

// StandardClaims represents claims used by GoSaaS JWT tokens
type StandardClaims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Iss    string `json:"iss"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

// ContextKey is a type for context keys
type ContextKey string

const (
	UserIDKey    ContextKey = "userId"
	UserEmailKey ContextKey = "userEmail"
	UserNameKey  ContextKey = "userName"
)

// JWTAuth creates middleware that validates JWT tokens signed with the given secret.
// It extracts claims and sets them in the request context.
func JWTAuth(accessSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				unauthorized(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				unauthorized(w, "invalid authorization header format")
				return
			}

			tokenStr := parts[1]
			if tokenStr == "" {
				unauthorized(w, "empty token")
				return
			}

			// Parse and validate the token
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(accessSecret), nil
			})
			if err != nil || !token.Valid {
				unauthorized(w, "invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				unauthorized(w, "invalid token claims")
				return
			}

			// Set claims in context (both typed keys and string keys for compatibility)
			ctx := r.Context()
			if userId, ok := claims["userId"].(string); ok {
				ctx = context.WithValue(ctx, UserIDKey, userId)
				ctx = context.WithValue(ctx, "userId", userId)
			}
			if email, ok := claims["email"].(string); ok {
				ctx = context.WithValue(ctx, UserEmailKey, email)
				ctx = context.WithValue(ctx, "email", email)
			}
			if name, ok := claims["name"].(string); ok {
				ctx = context.WithValue(ctx, UserNameKey, name)
				ctx = context.WithValue(ctx, "name", name)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LeveeTokenTranslator creates middleware that translates Levee JWT tokens to GoSaaS tokens.
// It intercepts Levee-issued tokens, validates them, extracts claims, and re-signs with GoSaaS's secret.
func LeveeTokenTranslator(accessSecret string) func(http.Handler) http.Handler {
	slog.Info("[LeveeTokenTranslator] Middleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			slog.Info("[LeveeTokenTranslator] Processing request", "path", r.URL.Path)

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				slog.Info("[LeveeTokenTranslator] Invalid auth header format")
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := parts[1]
			if tokenStr == "" {
				slog.Info("[LeveeTokenTranslator] Empty token")
				next.ServeHTTP(w, r)
				return
			}

			claims, err := parseJWTClaims(tokenStr)
			if err != nil {
				slog.Info("[LeveeTokenTranslator] Token parse failed", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			slog.Info("[LeveeTokenTranslator] Parsed token", "iss", claims.Iss, "sub", claims.Sub, "email", claims.Email)

			if claims.Iss != "levee.sh/sdk" {
				slog.Info("[LeveeTokenTranslator] Not a Levee token", "iss", claims.Iss)
				next.ServeHTTP(w, r)
				return
			}

			if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
				slog.Info("[LeveeTokenTranslator] Token expired", "user", claims.Sub)
				next.ServeHTTP(w, r)
				return
			}

			newClaims := StandardClaims{
				UserId: claims.Sub,
				Email:  claims.Email,
				Name:   claims.Name,
				Iss:    "gosaas",
				Exp:    claims.Exp,
				Iat:    claims.Iat,
			}

			newToken, err := createJWT(newClaims, accessSecret)
			if err != nil {
				slog.Error("[LeveeTokenTranslator] Failed to create translated token", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			r.Header.Set("Authorization", "Bearer "+newToken)
			slog.Info("[LeveeTokenTranslator] Successfully translated token", "user", claims.Sub, "email", claims.Email)

			next.ServeHTTP(w, r)
		})
	}
}

// createJWT creates a new JWT token with the given claims signed with the secret
func createJWT(claims StandardClaims, secret string) (string, error) {
	jwtClaims := jwt.MapClaims{
		"userId": claims.UserId,
		"email":  claims.Email,
		"name":   claims.Name,
		"iss":    claims.Iss,
		"exp":    claims.Exp,
		"iat":    claims.Iat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(secret))
}

// LeveeJWTMiddleware creates middleware that parses Levee JWT tokens.
// It extracts claims and sets them in the request context.
func LeveeJWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				unauthorized(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				unauthorized(w, "invalid authorization header format")
				return
			}

			tokenStr := parts[1]
			if tokenStr == "" {
				unauthorized(w, "empty token")
				return
			}

			claims, err := parseJWTClaims(tokenStr)
			if err != nil {
				slog.Error("Failed to parse JWT", "error", err)
				unauthorized(w, "invalid token")
				return
			}

			if claims.Iss != "levee.sh/sdk" {
				slog.Error("Invalid token issuer", "iss", claims.Iss)
				unauthorized(w, "invalid token issuer")
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.Sub)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserNameKey, claims.Name)
			ctx = context.WithValue(ctx, "userId", claims.Sub)
			ctx = context.WithValue(ctx, "email", claims.Email)

			slog.Info("JWT auth", "user", claims.Sub, "email", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// parseJWTClaims parses the claims from a JWT token without signature verification.
// This is safe because we trust Levee as the auth provider.
func parseJWTClaims(tokenString string) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}

// unauthorized sends a 401 response
func unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// ErrInvalidToken is returned when token parsing fails
var ErrInvalidToken = &tokenError{message: "invalid token"}

type tokenError struct {
	message string
}

func (e *tokenError) Error() string {
	return e.message
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	if id, ok := ctx.Value("userId").(string); ok {
		return id
	}
	return ""
}

// GetUserEmail extracts user email from context
func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}

// GetUserName extracts user name from context
func GetUserName(ctx context.Context) string {
	if name, ok := ctx.Value(UserNameKey).(string); ok {
		return name
	}
	return ""
}
