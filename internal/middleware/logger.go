package middleware

import (
	"time"
	"go.uber.org/zap"

	"go-backend/internal/logger"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		requestID, _ := c.Get("request_id")

		logger.Log.Info("HTTP Request",
			// structured logs 🔥
			zap.Any("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", duration),
		)
	}
}