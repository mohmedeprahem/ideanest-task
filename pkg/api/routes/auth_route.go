package route

import (
	"ideanest-task/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
    authController := controllers.NewAuthController()

    router.POST("/signup", authController.SignUp)
}
