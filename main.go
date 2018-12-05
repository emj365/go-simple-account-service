package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/emj365/account/lib"
	"github.com/emj365/account/models"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", postUsers).Methods("POST")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/auth", authUser).Methods("POST")
	fmt.Println("server is running on 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
	defer models.CloseDB()
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := &[]models.User{}
	models.GetDB().Find(&users)
	lib.Resonponse(w, http.StatusOK, users)
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	count := 0
	models.GetDB().Model(&models.User{}).Where("name = ?", user.Name).Count(&count)
	if count > 0 {
		lib.Resonponse(w, http.StatusConflict, map[string]interface{}{"name": user.Name})
		return
	}

	var err error
	user.Password, user.Salt, err = lib.GenPasswordHash(user.Password)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		lib.ResonponseServerError(w)
		return
	}

	models.GetDB().Create(user)
	lib.Resonponse(w, http.StatusCreated, user)
}

func authUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	foundUser := &models.User{Name: user.Name}
	models.GetDB().Find(foundUser)
	if !lib.Auth(foundUser, user.Name, user.Password) {
		w.WriteHeader(http.StatusForbidden)
	}

	jwt, err := lib.GetJWT(foundUser.ID)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		lib.ResonponseServerError(w)
		return
	}

	w.Header().Set("Authorization", "Bearer "+jwt)
	lib.Resonponse(w, http.StatusOK, map[string]interface{}{"jwt": jwt})
}
