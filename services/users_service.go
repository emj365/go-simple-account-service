package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/emj365/account/libs"
	"github.com/emj365/account/models"
	uuid "github.com/satori/go.uuid"
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

func CheckUserAlreadyExist(ch chan bool, name string, w http.ResponseWriter) {
	user := models.User{Name: name}
	exist := user.NameExistence()
	if exist {
		libs.Resonponse(w, http.StatusConflict, map[string]interface{}{"name": name})
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

	uuid, err := uuid.NewV4()

	if err != nil {
		log.Printf("Something went wrong: %s", err)
		libs.ResonponseServerError(w)
		ch <- false
		return
	}

	*salt = uuid.String()
	*hashedPassword = libs.HashPassword(password, *salt)

	ch <- true
}
