package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const deliveryTokenTTL = 5 * time.Minute

type DeliveryClaims struct {
	UserId     uint64 `json:"userId"`
	TradeNo    int64  `json:"tradeNo"`
	UniqueId   int32  `json:"uniqueId"`
	PayChannel string `json:"payChannel"`
	jwt.RegisteredClaims
}

func GenerateDeliveryJWT(userID uint64, tradeNo int64, uniqueId int32, payChannel string) (string, error) {
	claims := DeliveryClaims{
		UserId:     userID,
		UniqueId:   uniqueId,
		TradeNo:    tradeNo,
		PayChannel: payChannel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(deliveryTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateDeliveryJWT(tokenStr string) (*DeliveryClaims, error) {
	tkn, err := jwt.ParseWithClaims(tokenStr, &DeliveryClaims{},
		func(t *jwt.Token) (any, error) { return secretKey, nil })
	if err != nil || !tkn.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims := tkn.Claims.(*DeliveryClaims)

	return claims, nil
}
