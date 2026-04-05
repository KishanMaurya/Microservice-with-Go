package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-backend/internal/auth"
	"go-backend/internal/cache"
	"go-backend/internal/kafka"
	"go-backend/internal/models"
	"go-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req models.LoginRequest

	// 🔹 Step 1: Bind request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔹 Step 2: Get user from DB
	repo := &repositories.UserRepository{}
	user, err := repo.GetUserByEmail(req.Email)

	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 🔥 Step 3: Password validation
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 🔥 Step 4: Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access token error"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh token error"})
		return
	}

	// 🔥 Step 5: Store access token (optional - session tracking)
	err = cache.RDB.Set(cache.Ctx, accessToken, "active", 15*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
		return
	}

	// 🔥 Step 6: Store refresh token (IMPORTANT)
	err = cache.RDB.Set(cache.Ctx, refreshToken, user.ID, 7*24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
		return
	}

	// 🔥 Step 7: Multi-device session tracking
	sessionKey := "sessions:" + user.ID

	cache.RDB.SAdd(cache.Ctx, sessionKey, accessToken)
	cache.RDB.SAdd(cache.Ctx, sessionKey, refreshToken)
	cache.RDB.Expire(cache.Ctx, sessionKey, 7*24*time.Hour)

	// 🔥 Step 8: Publish login event to Kafka
	kafka.Publish("user.login", fmt.Sprintf("User %s logged in", user))

	// 🔹 Step 8: Response
	c.JSON(http.StatusOK, gin.H{
		"message":       "login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// 🔥 Check Redis
	userID, err := cache.RDB.Get(cache.Ctx, body.RefreshToken).Result()
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid refresh token"})
		return
	}

	// 🔥 Validate JWT
	token, err := auth.ValidateToken(body.RefreshToken)
	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	// 🔥 Generate new access token
	newAccessToken, err := auth.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "token error"})
		return
	}

	c.JSON(200, gin.H{
		"access_token": newAccessToken,
	})
}

func Logout(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format"})
		return
	}

	accessToken := parts[1]

	// 🔥 Step 1: blacklist access token (short-lived anyway)
	err := cache.RDB.Set(cache.Ctx, accessToken, "blacklisted", 15*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
		return
	}

	// 🔥 Step 2: delete refresh token (IMPORTANT)
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err == nil && body.RefreshToken != "" {
		cache.RDB.Del(cache.Ctx, body.RefreshToken)
	}

	// 🔥 Step 3: optional - remove from session set
	userID := c.GetString("user_id")
	sessionKey := "sessions:" + userID

	cache.RDB.SRem(cache.Ctx, sessionKey, accessToken)
	if body.RefreshToken != "" {
		cache.RDB.SRem(cache.Ctx, sessionKey, body.RefreshToken)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

func LogoutAll(c *gin.Context) {

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	key := "sessions:" + userID

	// 🔹 Step 1: get all tokens
	tokens, err := cache.RDB.SMembers(cache.Ctx, key).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
		return
	}

	// 🔹 Step 2: process tokens
	for _, token := range tokens {

		// 🔥 Try to parse token to detect type
		parsedToken, err := auth.ValidateToken(token)
		if err == nil && parsedToken != nil && parsedToken.Valid {

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if ok {
				tokenType, _ := claims["type"].(string)

				// 🔥 Access token → blacklist
				if tokenType == "access" {
					cache.RDB.Set(cache.Ctx, token, "blacklisted", 15*time.Minute)
				}

				// 🔥 Refresh token → delete
				if tokenType == "refresh" {
					cache.RDB.Del(cache.Ctx, token)
				}
			}
		} else {
			// fallback → just delete
			cache.RDB.Del(cache.Ctx, token)
		}
	}

	// 🔹 Step 3: delete session set
	cache.RDB.Del(cache.Ctx, key)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out from all devices",
	})
}