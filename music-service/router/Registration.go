package router

import (
	"context"
	"database/sql"
	"errors"
	"music-service/models"
	"music-service/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if repositories.SQLDB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db is not initialized"})
		return
	}
	userRepo := repositories.NewUserRepository(repositories.SQLDB)

	// username уникален
	if _, err := userRepo.GetByUsername(context.Background(), input.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// email уникален
	if _, err := userRepo.GetByEmail(context.Background(), input.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(passwordHash),
		Role:     "user",
	}

	if err := userRepo.Create(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}
