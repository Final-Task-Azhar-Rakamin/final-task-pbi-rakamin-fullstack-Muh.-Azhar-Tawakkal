package helpers

import (
	"finalapi/app"
	"finalapi/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user models.User) (string, error) {
	expTime := time.Now().Add(time.Minute * 3600)
	claims := &app.JWTClaim{
		Id:    int(user.Id),
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// medeklarasikan algoritma yang akan digunakan untuk signing
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signed token
	token, err := tokenAlgo.SignedString(app.JWT_KEY)
	return token, err
}

func ValidateUser(c *gin.Context) bool {
	id, _ := c.Get("id")
	dataId, err := strconv.ParseFloat(c.Param("userId"), 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return false
	}
	return id == dataId
}

func ValidateOwner(c *gin.Context, photo models.Photo) bool {
	id, _ := c.Get("id")
	photoId := photo.UserID
	return int64(id.(float64)) == photoId
}
