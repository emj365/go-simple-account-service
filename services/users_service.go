package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/emj365/go-simple-account-service/libs"
	"github.com/emj365/go-simple-account-service/models"
)

func GetUserFromRequest(
	w http.ResponseWriter, r *http.Request, user *models.User) bool {
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil || user.Name == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func GetUserID(r *http.Request) uint {
	return r.Context().Value(libs.ContextUserID).(uint)
}

func GetJWT(userID uint) (string, error) {
	now := time.Now()
	jsonWebToken := libs.JsonWebToken{SecretKey: libs.GetSecretKey(),
		Claims: jwt.StandardClaims{
			Subject:   strconv.Itoa(int(userID)),
			ExpiresAt: int64(now.Add(1 * time.Hour).Unix()),
		},
		Token: "",
	}

	err := jsonWebToken.GenToken()
	if err != nil {
		return "", err
	}

	return jsonWebToken.Token, nil
}
