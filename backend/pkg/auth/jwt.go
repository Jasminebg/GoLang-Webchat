package auth

import (
	"fmt"
	"time"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/models"
	"github.com/dgrijalva/jwt-go"
)

const hmacSecret = "135DSF4"
const defaultExpireTime = 604800

type Claims struct {
	Id       string `json:"id"`
	User     string `json:"user"`
	Color    string `json:"color"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func (c *Claims) GetId() string {
	return c.Id
}
func (c *Claims) GetName() string {
	return c.User
}
func (c *Claims) GetColor() string {
	return c.Color
}
func (c *Claims) GetPassword() string {
	return c.Password
}

func CreateJWTToken(user models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        user.GetId(),
		"user":      user.GetName(),
		"color":     user.GetColor(),
		"ExpiresAt": time.Now().Unix() + defaultExpireTime,
	})

	tokenString, err := token.SignedString([]byte(hmacSecret))

	return tokenString, err
}

func ValidateToken(tokenString string) (models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method:%v", token.Header["alg"])
		}

		return []byte(hmacSecret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
