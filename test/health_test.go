package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Success(t *testing.T) {

	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create router
	router := gin.Default()
	router.GET("/health", handlers.HealthCheck)

	// Create request
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"UP"}`, w.Body.String())
}