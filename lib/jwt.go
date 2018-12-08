package lib

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetJWT return jwt with userId
func GetJWT(userID uint) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": int(userID),
		"exp": now.Add(1 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("hmacSampleSecret"))
}

func DecodeJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("hmacSampleSecret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
