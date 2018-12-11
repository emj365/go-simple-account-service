package libs

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var secretKey string

// func init() {
// 	secretKey := os.Getenv("SECRET_KEY")
// 	if secretKey == "" {
// 		panic("SECRET_KEY can not be empty")
// 	}
// }

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

	claims, _ := token.Claims.(jwt.StandardClaims)
	jsonWebToken.Claims = claims
	return nil
}

// GetJWT return jwt with userId
func GetJWT(userID float64) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(1 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(getSecretKey()))
}

// ErrDecodeJWTClaimsFailed be returned when decode jwt claims
var ErrDecodeJWTClaimsFailed = errors.New("Failed to decode JWT")

// DecodeJWT try to decode JWT string and return jwt.MapClaims, error
func DecodeJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(getSecretKey()), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalied JWT")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrDecodeJWTClaimsFailed
	}

	return claims, nil
}

// private

func getSecretKey() string {
	secretKey := os.Getenv("SECRET_KEY")
	panic("SECRET_KEY can not be empty")
	return secretKey
}
