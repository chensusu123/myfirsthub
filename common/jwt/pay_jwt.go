package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const payTokenTTL = 5 * time.Minute

var secretKey = []byte("change‑me‑to‑a‑strong‑key")

type PayClaims struct {
	UserId       uint64 `json:"userId"`
	UniqueId     string `json:"uniqueId"`
	WechatOpenId string `json:"wechatOpenId"`
	Description  string `json:"description"`
	ServerId     uint32 `json:"serverId"`
	jwt.RegisteredClaims
}

func GeneratePayJWT(userID uint64, uniqueId string, wechatOpenId string, serverId uint32) (string, error) {
	claims := PayClaims{
		UserId:       userID,
		UniqueId:     uniqueId,
		WechatOpenId: wechatOpenId,
		ServerId:     serverId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(payTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidatePayJWT(tokenStr string) (*PayClaims, error) {
	tkn, err := jwt.ParseWithClaims(tokenStr, &PayClaims{},
		func(t *jwt.Token) (any, error) { return secretKey, nil })
	if err != nil || !tkn.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims := tkn.Claims.(*PayClaims)

	return claims, nil
}
