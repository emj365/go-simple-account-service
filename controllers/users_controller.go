package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/emj365/account/lib"
	"github.com/emj365/account/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "getUsers")

	users := models.FindAllUsers()
	lib.Resonponse(w, http.StatusOK, users)
}

func PostUsers(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "postUsers")

	user := models.User{}
	if !lib.GetUserFromRequest(w, r, &user) {
		return
	}

	var hashedPassword, salt string
	ch := make(chan bool)
	go lib.CheckUserAlreadyExist(ch, user.Name, w)
	go lib.GenHashedPassword(ch, user.Password, &hashedPassword, &salt, w)

	countOfRecived := 0
	for countOfRecived < 2 {
		select {
		case successed := <-ch:
			if !successed {
				return
			}

			countOfRecived++
		}
	}

	createdUser := models.User{Name: user.Name, Password: hashedPassword, Salt: salt}
	models.CreateUser(&createdUser)
	lib.Resonponse(w, http.StatusCreated, createdUser)
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "auth")

	user := models.User{}
	if !lib.GetUserFromRequest(w, r, &user) {
		return
	}

	foundUser := models.GetUserForAuth(user.Name)
	if foundUser.Name == "" || !lib.Auth(foundUser, user.Name, user.Password) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	jwt, err := lib.GetJWT(foundUser.ID)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		lib.ResonponseServerError(w)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)
	lib.Resonponse(w, http.StatusOK, map[string]interface{}{"jwt": jwt})
}
