package app

import (
	"github.com/gin-gonic/gin"

	route "ideanest-task/pkg/api/routes"
	database "ideanest-task/pkg/database/mongodb"
)

func RunApp() {
	// Initialize MongoDB connection
	database.DB = database.ConnectDB()

	// Set up Gin router
	router := gin.Default()

	// Set up routes
	route.AuthRoute(router)

	// Start server
	router.Run(":8080")
}
