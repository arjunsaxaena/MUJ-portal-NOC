package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsHandlerMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},                   // List of allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // List of allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // List of allowed headers
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},          // Headers that can be exposed
		AllowCredentials: true,                                                // Allow cookies
		MaxAge:           12 * time.Hour,                                      // Cache preflight request duration
	})
}
