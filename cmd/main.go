package main

import (
	app "ideanest-task/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	app.Run()
	r := gin.Default()

	r.Run()
}
