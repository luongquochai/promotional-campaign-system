package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/database/redis"
	"github.com/luongquochai/promotional-campaign-system/models"
	"github.com/luongquochai/promotional-campaign-system/services"
	"golang.org/x/crypto/bcrypt"
)

// Register a new user
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /user/register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := services.RegisterService(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Authenticate user with email and password
// @Summary Authenticate user
// @Description Authenticate user and return JWT token
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.User true "User credentials"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /user/login [post]
func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.LoginService(c, input.Email, input.Password)
	if err != nil {
		return
	}

	// Generate JWT
	token, err := config.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token"})
		return
	}

	// Store user session in Redis
	err = redis.Set(token, strconv.FormatUint(uint64(user.ID), 10), 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session in Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
