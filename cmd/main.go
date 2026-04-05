package main

import (
	"go-backend/internal/cache"
	"go-backend/internal/database"
	"go-backend/internal/kafka"
	"go-backend/internal/logger"
	"go-backend/internal/middleware"
	"go-backend/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger.Init()
	defer logger.Sync()

	database.ConnectMongo()
	cache.InitRedis()
	kafka.InitProducer()

	router := gin.New()

	router.Use(middleware.RequestID())
	router.Use(middleware.Logger()) // custom logger
	router.Use(gin.Recovery())

	routes.RegisterRoutes(router)

	router.Run(":8080")
}
