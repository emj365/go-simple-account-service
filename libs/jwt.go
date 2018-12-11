package libs

import (
	"errors"
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

type JsonWebToken struct {
	SecretKey string
	Token     string
	Claims    jwt.StandardClaims
}

func (jsonWebToken *JsonWebToken) GenToken() error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jsonWebToken.Claims)

	var err error
	jsonWebToken.Token, err = token.SignedString([]byte(jsonWebToken.SecretKey))
	return err
}

func (jsonWebToken *JsonWebToken) Decode() error {
	token, err := jwt.Parse(jsonWebToken.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jsonWebToken.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalied JWT")
	}

	var claims jwt.MapClaims
	var ok bool
	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("Invalied JWT Claims")
	}

	var sub string
	sub, ok = claims["sub"].(string)
	if !ok {
		return errors.New("Invalied JWT Claims")
	}

	var exp float64
	if exp, ok = claims["exp"].(float64); !ok {
		return errors.New("Invalied JWT Claims")
	}

	jsonWebToken.Claims.Subject = sub
	jsonWebToken.Claims.ExpiresAt = int64(exp)

	return nil
}

func GetSecretKey() string {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic("SECRET_KEY can not be empty")
	}
	return secretKey
}
