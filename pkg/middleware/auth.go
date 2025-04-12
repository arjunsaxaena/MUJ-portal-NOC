package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	blacklistedTokens = make(map[string]bool)
	blacklistMutex    sync.RWMutex
)

func BlacklistToken(token string) {
	token = strings.TrimPrefix(token, "Bearer ")

	blacklistMutex.Lock()
	blacklistedTokens[token] = true
	blacklistMutex.Unlock()
}

func IsTokenBlacklisted(token string) bool {
	token = strings.TrimPrefix(token, "Bearer ")

	blacklistMutex.RLock()
	isBlacklisted := blacklistedTokens[token]
	blacklistMutex.RUnlock()

	return isBlacklisted
}

func AuthMiddleware(jwtSecretKey string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if IsTokenBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		userRole := claims["role"].(string)

		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Set("id", claims["id"])
		c.Set("email", claims["email"])
		c.Set("role", userRole)
		c.Set("department", claims["department"])

		if roleType, exists := claims["roleType"]; exists {
			c.Set("roleType", roleType)
		}

		c.Next()
	}
}
