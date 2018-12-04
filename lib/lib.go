package lib

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// GetJWT return jwt with userId
func GetJWT(userId int) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userId),
		"exp": now.Add(1 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("hmacSampleSecret"))
}

func GenPasswordHash(password string) (string, string, error) {
	salt, err := uuid.NewV4()
	if err != nil {
		return "", "", err
	}

	return HashPassword(password, salt.String()), salt.String(), nil
}

func HashPassword(rawPassword string, salt string) string {
	hasher := md5.New()
	hasher.Write([]byte(rawPassword + salt))
	password := hex.EncodeToString(hasher.Sum(nil))
	return password
}
