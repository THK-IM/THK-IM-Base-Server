package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserTokenClaims struct {
	Id        int64
	Username  *string
	ExpiresAt int64
	Issuer    string
}

func GenerateUserToken(id int64, serverName, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"expiresAt": expireTime.Unix(),
		"issuer":    serverName,
	})
	hmacSampleSecret := []byte(password)
	return token.SignedString(hmacSampleSecret)
}

func CheckUserToken(token, password string) (*int64, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(password), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		id, okId := claims["id"].(float64)
		if okId {
			idInt64 := int64(id)
			return &idInt64, nil
		}
	}

	return nil, errors.New("invalid token")

}
