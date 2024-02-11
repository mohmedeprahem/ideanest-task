package controllers

import (
	"ideanest-task/pkg/database/mongodb/models"
	"ideanest-task/pkg/database/mongodb/repository"
	"ideanest-task/pkg/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userRepository *repository.UserRepository
}

var validate = validator.New()

func NewAuthController() *UserController {
	return &UserController{
		userRepository: repository.NewUserRepository(),
	}
}

func (c *UserController) Create(ctx *gin.Context) {
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
