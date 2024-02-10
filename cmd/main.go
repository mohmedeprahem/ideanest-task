package main

import (
	database "ideanest-task/pkg/database/mongodb"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()
	r := gin.Default()

	r.Run()
}
