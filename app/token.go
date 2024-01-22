package app

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Email string
	jwt.RegisteredClaims
}

func verifyToken(tok string, key []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tok, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		fmt.Println("verifyToken uses key", string(key))
		return key, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	cl, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", errors.New("unable to extract user claims")
	}
	return cl.Email, nil
}
