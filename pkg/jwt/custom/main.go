package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func main() {
	type CustomClaims struct {
		// 额外的信息
		Exp string `json:"exp"`
		jwt.RegisteredClaims
	}
	claims := CustomClaims{
		Exp: "Foo",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Owner",
			Subject:   "JSON Web Token",
			Audience:  []string{"Someone_else"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "123",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// salt
	secret := []byte("secret")
	sign, err := token.SignedString(secret)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token.Method)
	fmt.Println(sign)

	// 解码
	parseToken, err := jwt.ParseWithClaims(sign, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if claims, ok := parseToken.Claims.(*CustomClaims); ok && parseToken.Valid {
		fmt.Println(claims.Exp, claims.ID, claims.Issuer)
	} else {
		fmt.Println(err)
	}
}
