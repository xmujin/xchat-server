package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var mySigningKey = []byte("我真的服了")

func GenerateToken(username string) string {
	claims := MyCustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "xchat-server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	return ss
}

func ParseToken(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, errors.New("unknown claims type, cannot proceed")
	}
	return claims, nil
}
