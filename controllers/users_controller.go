package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/emj365/go-simple-account-service/libs"
	"github.com/emj365/go-simple-account-service/models"
	"github.com/emj365/go-simple-account-service/services"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	defer libs.TimeTrack(time.Now(), "getUsers")

	users := models.GetAllUser()
	libs.Resonponse(w, http.StatusOK, users)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	defer libs.TimeTrack(time.Now(), "getMe")

	userID := r.Context().Value("userID")

	user := models.User{}
	models.FindUserByID(&user, uint(userID.(float64)))
	libs.Resonponse(w, http.StatusOK, user)
}

func PostUsers(w http.ResponseWriter, r *http.Request) {
	defer libs.TimeTrack(time.Now(), "postUsers")

	user := models.User{}
	if !services.GetUserFromRequest(w, r, &user) {
		return
	}

	exist := user.NameExistence()
	if exist {
		libs.Resonponse(w, http.StatusConflict, map[string]interface{}{"name": user.Name})
		return
	}

	err := user.Create()
	if err != nil {
		log.Printf("error: %v\n", err)
		libs.ResonponseServerError(w)
	}

	libs.Resonponse(w, http.StatusCreated, user)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	defer libs.TimeTrack(time.Now(), "auth")

	user := models.User{}
	if !services.GetUserFromRequest(w, r, &user) {
		return
	}

	if !user.Auth() {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	jwt, err := libs.GetJWT(float64(user.ID))
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		libs.ResonponseServerError(w)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)
	libs.Resonponse(w, http.StatusOK, map[string]interface{}{"jwt": jwt})
}

func JWT(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	jwt, error := libs.GetJWT(userID.(float64))
	if error != nil {
		libs.ResonponseServerError(w)
		log.Println("GetJWT Error")
		return
	}

	libs.Resonponse(w, http.StatusOK, map[string]string{"jwt": jwt})
}
