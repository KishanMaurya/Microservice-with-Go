package handlers

import (
	"fmt"
	"net/http"
	"time"

	// "go-backend/internal/logger"
	"go-backend/internal/kafka"
	"go-backend/internal/models"
	"go-backend/internal/repositories"

	// "go-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	loggerInterface, _ := c.Get("logger")
	log := loggerInterface.(*zap.Logger)

	log.Info("User created",
		zap.String("email", req.Email),
	)
	// Bind request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Info("Invalid input detected")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	user := models.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo := &repositories.UserRepository{}
	err = repo.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kafka.Publish("user.register", fmt.Sprintf("User %s Registered in", user))
	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
	})
}
