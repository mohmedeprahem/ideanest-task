package route

import (
	"ideanest-task/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
    authController := controllers.NewAuthController()

    router.POST("/signup", authController.SignUp)
		router.POST("/signin", authController.SignIn)
		router.POST("/refresh-token", authController.RefreshToken)
}

