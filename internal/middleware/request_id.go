package middleware

import (
	"go-backend/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// create child logger with request_id
		reqLogger := logger.Log.With(
			zap.String("request_id", requestID),
		)

		// store in context
		c.Set("request_id", requestID)
		c.Set("logger", reqLogger)

		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}