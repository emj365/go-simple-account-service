package lib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/emj365/account/models"
	uuid "github.com/satori/go.uuid"
)

// GetJWT return jwt with userId
func GetJWT(userID uint) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(userID)),
		"exp": now.Add(1 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("hmacSampleSecret"))
}

func GenPasswordHash(password string) (string, string, error) {
	salt, err := uuid.NewV4()
	if err != nil {
		return "", "", err
	}

	return hashPassword(password, salt.String()), salt.String(), nil
}

func Auth(user models.User, name string, password string) bool {
	hash := hashPassword(password, user.Salt)
	if name == user.Name && hash == user.Password {
		return true
	}

	return false
}
func GetUserFromRequest(
	w http.ResponseWriter, r *http.Request, user *models.User) bool {
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil || user.Name == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func CheckUserAlreadyExist(ch chan bool, name string, w http.ResponseWriter) {
	exist := models.CheckUserAlreadyExist(name)
	if exist {
		Resonponse(w, http.StatusConflict, map[string]interface{}{"name": name})
		ch <- false
		return
	}

	ch <- true
}

func GenHashedPassword(
	ch chan bool,
	password string,
	hashedPassword *string,
	salt *string,
	w http.ResponseWriter) {
	var err error
	*hashedPassword, *salt, err = GenPasswordHash(password)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		ResonponseServerError(w)
		ch <- false
		return
	}

	ch <- true
}

// private

func hashPassword(rawPassword string, salt string) string {
	hasher := md5.New()
	hasher.Write([]byte(rawPassword + salt))
	password := hex.EncodeToString(hasher.Sum(nil))
	return password
}
