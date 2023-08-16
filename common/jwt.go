package common

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 签名密钥
var jwtSecretKey = []byte("hi")

func GenerateToken(userId uint, username string) string {

	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour).Unix()

	log.Println("expireTime:", expireTime)

	claims := Claims{
		ID:       userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			IssuedAt:  nowTime.Unix(),
			Issuer:    "DouShen",
		},
	}
	// 使用用于签名的算法和令牌
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 创建JWT字符串
	if token, err := tokenClaims.SignedString(jwtSecretKey); err != nil {
		log.Println("generate token fail!")
		return "fail"
	} else {
		log.Println("generate token success!")
		return token
	}
}

// ParseToken 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
