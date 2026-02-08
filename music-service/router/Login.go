package router

import (
	"context"
	"database/sql"
	"errors"
	"music-service/models"
	"music-service/repositories"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if repositories.SQLDB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db is not initialized"})
		return
	}
	userRepo := repositories.NewUserRepository(repositories.SQLDB)

	// сначала пробуем как username
	userFound, err := userRepo.GetByUsername(context.Background(), input.Login)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		// если username не найден — пробуем как email
		userFound, err = userRepo.GetByEmail(context.Background(), input.Login)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server misconfigured"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	token, err := generateToken.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
