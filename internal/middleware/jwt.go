package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const ContextUserIDKey = "userId"

type JWTClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// if c.Request().URL.Path == "/api/v1/login" || c.Request().URL.Path == "/api/v1/auth/login" || c.Request().URL.Path == "/api/v1/auth/signup" {
			// 	// if c.Request().Method == http.MethodOptions {
			// 	return next(c)
			// }

			// authHeader := c.Request().Header.Get("Authorization")
			// if authHeader == "" {
			// 	return c.JSON(http.StatusUnauthorized, map[string]any{
			// 		"error": map[string]any{
			// 			"message": "missing Authorization header",
			// 		},
			// 	})
			// }

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"error": map[string]any{
						"message": "missing Authorization header",
					},
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"error": map[string]any{
						"message": "invalid Authorization header format",
					},
				})
			}

			tokenStr := parts[1]
			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"error": map[string]any{
						"message": "invalid or expired token",
					},
				})
			}

			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"error": map[string]any{
						"message": "token expired",
					},
				})
			}

			c.Set(ContextUserIDKey, claims.UserID)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) string {
	if v := c.Get(ContextUserIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// package middleware

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/labstack/echo/v4"
// )

// const ContextUserIDKey = "userId"

// type JWTClaims struct {
// 	UserID string `json:"userId"`
// 	jwt.RegisteredClaims
// }

// func JWTAuth(secret string) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			// ✅ Allow preflight requests through
// 			if c.Request().Method == http.MethodOptions {
// 				return next(c)
// 			}

// 			authHeader := c.Request().Header.Get("Authorization")
// 			if authHeader == "" {
// 				return c.JSON(http.StatusUnauthorized, map[string]any{
// 					"error": map[string]any{
// 						"message": "missing Authorization header",
// 					},
// 				})
// 			}

// 			parts := strings.SplitN(authHeader, " ", 2)
// 			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
// 				return c.JSON(http.StatusUnauthorized, map[string]any{
// 					"error": map[string]any{
// 						"message": "invalid Authorization header format",
// 					},
// 				})
// 			}

// 			tokenStr := parts[1]
// 			claims := &JWTClaims{}
// 			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
// 				return []byte(secret), nil
// 			})
// 			if err != nil || !token.Valid {
// 				return c.JSON(http.StatusUnauthorized, map[string]any{
// 					"error": map[string]any{
// 						"message": "invalid or expired token",
// 					},
// 				})
// 			}

// 			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
// 				return c.JSON(http.StatusUnauthorized, map[string]any{
// 					"error": map[string]any{
// 						"message": "token expired",
// 					},
// 				})
// 			}

// 			c.Set(ContextUserIDKey, claims.UserID)
// 			return next(c)
// 		}
// 	}
// }

// func GetUserID(c echo.Context) string {
// 	if v := c.Get(ContextUserIDKey); v != nil {
// 		if s, ok := v.(string); ok {
// 			return s
// 		}
// 	}
// 	return ""
// }
