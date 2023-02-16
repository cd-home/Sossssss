package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func main() {
	claims := &jwt.RegisteredClaims{
		Issuer:    "Owner",
		Subject:   "JSON Web Token",
		Audience:  []string{"Someone_else"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "123",
	}
	salt := []byte("salt")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, _ := token.SignedString(salt)
	fmt.Println(sign)

	// 解码
	parseToken, err := jwt.Parse(sign, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return salt, nil
	})
	if claims, ok := parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {
		fmt.Println(claims["jti"])
	} else {
		fmt.Println(err)
	}
}
