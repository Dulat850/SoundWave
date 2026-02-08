package router

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music-service/repositories"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
		return
	}

	authToken := strings.SplitN(authHeader, " ", 2)
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		return
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server misconfigured"})
		return
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	expAny, ok := claims["exp"]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has no exp"})
		return
	}

	var expUnix int64
	switch v := expAny.(type) {
	case float64:
		expUnix = int64(v)
	case int64:
		expUnix = v
	case int:
		expUnix = int64(v)
	default:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid exp type"})
		return
	}

	if time.Now().Unix() > expUnix {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return
	}

	idAny, ok := claims["id"]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has no id"})
		return
	}

	var userID int
	switch v := idAny.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	case int64:
		userID = int(v)
	default:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid id type"})
		return
	}

	if repositories.SQLDB == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db is not initialized"})
		return
	}
	userRepo := repositories.NewUserRepository(repositories.SQLDB)

	user, err := userRepo.GetByID(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.Set("currentUser", user)
	c.Next()
}
