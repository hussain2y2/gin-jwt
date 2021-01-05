package main

import (
	"go-jwt/auth"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	dsn := os.Getenv("REDIS_DSN")
	client = redis.NewClient(&redis.Options{Addr: dsn})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}
}

// CreateAuth ...
func CreateAuth(userID uint64, token *auth.Token) error {
	accessToken := time.Unix(token.AccessToken.Expires, 0)
	refreshToken := time.Unix(token.RefreshToken.Expires, 0)
	now := time.Now()

	errAccess := client.Set(token.AccessToken.UUID, strconv.Itoa(int(userID)), accessToken.Sub(now)).Err()

	if errAccess != nil {
		return errAccess
	}

	errRefresh := client.Set(token.RefreshToken.UUID, strconv.Itoa(int(userID)), refreshToken.Sub(now)).Err()

	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func main() {
	router := gin.Default()
	router.POST("/login", auth.Login)
	router.Run(":8080")
}
