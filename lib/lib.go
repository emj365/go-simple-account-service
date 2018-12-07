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
func GetJWT(userId uint) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(userId)),
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

func Auth(user models.User, name string, password string) bool {
	hash := HashPassword(password, user.Salt)
	if name == user.Name && hash == user.Password {
		return true
	}

	return false
}

func ResonponseServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": http.StatusText(http.StatusInternalServerError)})
}

func Resonponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
