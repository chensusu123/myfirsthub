package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 定义自定义Claims结构
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var secretKey = []byte("zU6W/(Y%,KX?-@q4m~tLy1_uhcekTNQg")

func GenerateToken(userID, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "your-app-name",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func main() {
	userID := "12345"
	username := "john_doe"
	token, _ := GenerateToken(userID, username)
	println(token)
	println("========================")
	println(validateToken(token))

	tokenv := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDUiLCJ1c2VybmFtZSI6ImpvaG5fZG9lIiwiZXhwIjoxNzQ3NTUyMjcwLCJpc3MiOiJ5b3VyLWFwcC1uYW1lIn0.91_jlpnGuJYvpJVbxebC5rwOhy0JFhyY-g2kZpuiLaY"
	println(validateToken(tokenv))
	// claims, _ := ValidateToken(token)
	// println(claims.UserID)
	// println(claims.Username)
}

type MyCustomClaims struct {
	jwt.StandardClaims
}

func VerifyWithCustomClaims(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func validateToken(tokenIn string) bool {
	_, err := VerifyWithCustomClaims(tokenIn)
	return err == nil
}
