package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	database "github.com/luongquochai/promotional-campaign-system/database/mysql"
	"github.com/luongquochai/promotional-campaign-system/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterService(c *gin.Context, user models.User) error {
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
		return err
	}

	return nil
}

func LoginService(c *gin.Context, email, password string) (models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return user, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return user, err
	}

	return user, nil
}
