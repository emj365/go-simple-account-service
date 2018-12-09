package lib

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func getSecretKey() string {
	secretKey := os.Getenv("SECRET_KEY")
	return secretKey
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

// DecodeJWT try to decode JWT string and return jwt.MapClaims, error
func DecodeJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(getSecretKey()), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); err == nil && ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
