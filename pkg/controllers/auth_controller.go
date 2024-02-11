package controllers

import (
	"ideanest-task/pkg/database/mongodb/models"
	"ideanest-task/pkg/database/mongodb/repository"
	"ideanest-task/pkg/utils"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userRepository *repository.UserRepository
}

type SignInRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}


var validate = validator.New()

func NewAuthController() *UserController {
	return &UserController{
		userRepository: repository.NewUserRepository(),
	}
}

func (c *UserController) SignUp(ctx *gin.Context) {
	// Get user data from request body
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user data
	if err := validate.Struct(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	if _, err := c.userRepository.FindByEmail(user.Email); err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = hashedPassword

	// Create user
	user.ID = primitive.NewObjectID()
	if err := c.userRepository.Create(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func (c *UserController) SignIn(ctx *gin.Context) {
	// Get user data from request body
	var input  SignInRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	user, err := c.userRepository.FindByEmail(input.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}


	// User exists, check password
	if err := utils.CheckPassword(user.Password, input.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}


	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}


	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Save refresh token in redis
	err = utils.SetRedisValue(refreshToken, true, 7 * 24 * time.Hour)
	if err != nil {
		log.Println(err)
	    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	    return
	}

	// Return access token and refresh token
	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

// refresh token
func (c *UserController) RefreshToken(ctx *gin.Context) {
	configApp, _ := utils.ReadAppConfig()
	// Get refresh token from request body
	var input RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if refresh token is valid
	if isNotValid, err := utils.IsTokenInvalid(input.RefreshToken, configApp.Jwt.RtSecret); isNotValid || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "12invalid refresh token"})
		return
	}

	// Check if refresh token exists
	IsActive, err := utils.GetRedisValue(input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		return
	}

	if IsActive != "1" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "is not active refresh token"})
		return
	}

	// revoke refresh token
	err = utils.DeleteRedisValue(input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	err = utils.SetRedisValue(input.RefreshToken, false, time.Hour*24*7)
	if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	// Save new refresh token in redis
	err = utils.SetRedisValue(refreshToken, true, time.Hour*24*7)
	if err != nil {
	    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	    return
	}

	// Return successful response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "refresh token successful",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})

}

// Revoke token
func (c *UserController) RevokeToken(ctx *gin.Context) {
	configApp, _ := utils.ReadAppConfig()
	// Get refresh token from request body
	var input RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if refresh token is valid
	if isNotValid, err := utils.IsTokenInvalid(input.RefreshToken, configApp.Jwt.RtSecret); isNotValid || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		return
	}

	// Check if refresh token exists
	value, err := utils.GetRedisValue(input.RefreshToken);
	if err != nil {
	    ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
	    return
	}

	if value != "1" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	// revoke refresh token
	err = utils.DeleteRedisValue(input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		return
	}

	err = utils.SetRedisValue(input.RefreshToken, false, time.Hour*24*7)
	if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "internal server error"})
			return
	}

	// Delete access token from Header if exists
	if ctx.GetHeader("Authorization") != "" {
		ctx.Request.Header.Del("Authorization")
	}

	// Return successful response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "revoke token successful",
	})

}
