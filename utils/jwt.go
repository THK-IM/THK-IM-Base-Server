package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
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
	strId := strconv.FormatInt(id, 36)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        strId,
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
		strId, okId := claims["id"].(string)
		if okId {
			idInt64, errParse := strconv.ParseInt(strId, 36, 64)
			return &idInt64, errParse
		}
	}

	return nil, errors.New("invalid token")

}
