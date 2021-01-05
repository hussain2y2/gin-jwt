package auth

import (
	"go-jwt/models"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/twinj/uuid"
)

var client *redis.Client

// TODO Replace with DB Record
var user = models.User{
	ID:       1,
	Email:    "engrhussainahmad@gmail.com",
	Password: "password",
}

// Token ...
type Token struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

// AccessToken ...
type AccessToken struct {
	Token   string
	UUID    string
	Expires int64
}

// RefreshToken ...
type RefreshToken struct {
	Token   string
	UUID    string
	Expires int64
}

func init() {
	dsn := os.Getenv("REDIS_DSN")
	client = redis.NewClient(&redis.Options{Addr: dsn})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}
}

// Login ...
func Login(c *gin.Context) {
	var usr models.User

	if err := c.ShouldBindJSON(&usr); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided!")
		return
	}

	if user.Email != usr.Email || user.Password != usr.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details!")
		return
	}

	token, err := CreateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	saveToken := CreateAuth(user.ID, token)

	if saveToken != nil {
		c.JSON(http.StatusUnprocessableEntity, saveToken.Error())
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken.Token, "refresh_token": token.RefreshToken.Token})

}

// CreateToken ...
func CreateToken(userID uint64) (*Token, error) {
	token := &Token{}

	token.AccessToken.Expires = time.Now().Add(time.Minute * 15).Unix()
	token.AccessToken.UUID = uuid.NewV4().String()

	token.RefreshToken.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	token.RefreshToken.UUID = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = token.AccessToken.UUID
	atClaims["userId"] = userID
	atClaims["expires"] = token.AccessToken.Expires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token.AccessToken.Token, err = at.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}

	rtClaims["refresh_uuid"] = token.RefreshToken.UUID
	rtClaims["userId"] = userID
	rtClaims["expires"] = token.RefreshToken.Expires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken.Token, err = rt.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))

	if err != nil {
		return nil, err
	}

	return token, nil
}

// CreateAuth ...
func CreateAuth(userID uint64, token *Token) error {
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
