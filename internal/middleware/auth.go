package middleware

import (
	"net/http"
	"strings"

	"go-backend/internal/auth"
	"go-backend/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 🔥 STEP 1: Check Redis (blacklist + session)
		val, err := cache.RDB.Get(cache.Ctx, tokenString).Result()

		if err != nil {
			// token not found in Redis → invalid session
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired or invalid"})
			c.Abort()
			return
		}

		if val == "blacklisted" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			c.Abort()
			return
		}

		// 🔥 STEP 2: Validate JWT
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// 🔥 STEP 3: Extract claims (JWT v5)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			c.Abort()
			return
		}

		// 🔥 STEP 4: Extract user_id safely
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
			c.Abort()
			return
		}

		// attach to context
		c.Set("user_id", userID)

		c.Next()
	}
}