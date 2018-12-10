package libs

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

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
	return secretKey
}
