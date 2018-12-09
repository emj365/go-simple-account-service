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

	users := models.GetAllUser()
	lib.Resonponse(w, http.StatusOK, users)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "getMe")

	userID := r.Context().Value("userID")

	user := models.User{}
	models.FindUserByID(&user, uint(userID.(float64)))
	lib.Resonponse(w, http.StatusOK, user)
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

	newUser := models.User{Name: user.Name, Password: hashedPassword, Salt: salt}
	err := newUser.Create()
	if err != nil {
		log.Printf("error: %v\n", err)
		lib.ResonponseServerError(w)
	}

	lib.Resonponse(w, http.StatusCreated, newUser)
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "auth")

	user := models.User{}
	if !lib.GetUserFromRequest(w, r, &user) {
		return
	}

	name, password := user.Name, user.Password

	found := user.FindForAuth()
	if !found || !lib.Auth(user, name, password) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	jwt, err := lib.GetJWT(float64(user.ID))
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		lib.ResonponseServerError(w)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)
	lib.Resonponse(w, http.StatusOK, map[string]interface{}{"jwt": jwt})
}

func JWT(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	jwt, error := lib.GetJWT(userID.(float64))
	if error != nil {
		lib.ResonponseServerError(w)
		log.Println("GetJWT Error")
		return
	}

	lib.Resonponse(w, http.StatusOK, map[string]string{"jwt": jwt})
}
