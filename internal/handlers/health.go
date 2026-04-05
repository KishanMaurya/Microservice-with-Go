package handlers

import (
	"go-backend/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	logger.Log.Info("Health check endpoint hit")
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}