package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	route "ideanest-task/pkg/api/routes"
	database "ideanest-task/pkg/database/mongodb"
	"ideanest-task/pkg/utils"
)


func RunApp(ctx context.Context) {

	// Initialize Redis connection
	redisClient := utils.NewRedisClient()

	// Initialize MongoDB connection
	database.DB = database.ConnectDB()

	// Ping Redis
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		fmt.Println(err)
	}

	// Set up Gin router
	router := gin.Default()

	// Set up routes
	route.AuthRoute(router)

	// Start server
	router.Run(":8080")
}
