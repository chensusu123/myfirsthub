package userprofile

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

func CreateAvatarToken(logger fklog.FKLogI) (string, error) {
	token, err := GenerateAvatarToken()
	if err != nil {
		logger.ErrorWF("CreateAvatarToken failed", zap.Error(err))
		return "", err
	}
	return token, nil
}

type CustomClaims struct {
	jwt.StandardClaims
}

const (
	validUntil = 3600 // 有效期
)

var secretKey = "hfV98QwSEPxcrfW2BpstGA=="

// GenerateAvatarToken 生成头像jwt token
func GenerateAvatarToken() (string, error) {
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(validUntil * time.Second).Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}

// VerifyWithCustomClaims 验证token
func VerifyWithCustomClaims(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	// todo 先简单实现，后期使用统一解决方案
	// 验证Token有效性并提取Claims
	// if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
	// 	if claims.Issuer != "user_avatar" {
	// 		return nil, fmt.Errorf("invalid issuer")
	// 	}
	// }
	return nil, err
}

func ValidateAvatarToken(tokenIn string) bool {
	_, err := VerifyWithCustomClaims(tokenIn)
	return err == nil
}
