package common

import (
	"building-distributed-app-in-gin-chapter06/api/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("a_secret_cre")

type Claims struct {
	UserName string
	jwt.StandardClaims
}

func ReleaseToken(form models.Form) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserName: form.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //token失效时间
			IssuedAt:  time.Now().Unix(),     //token发放时间
			Issuer:    "oceanlearn.tech",     //token发放者
			Subject:   "user token",          //主题
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	//token 成功生成
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, err
	})
	return token, claims, err
}
