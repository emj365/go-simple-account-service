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

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

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
	models.GetDB().Find(users)
	json.NewEncoder(w).Encode(users)
	return
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)
	var err error
	user.Password, user.Salt, err = lib.GenPasswordHash(user.Password)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}

	models.GetDB().Create(user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func authUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	if auth(user.Name, user.Password) {
		jwt, err := lib.GetJWT(123)
		if err != nil {
			fmt.Printf("Something went wrong: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Message(false, "can not generate JWT"))
			return
		}

		w.Header().Set("Authorization", "Bearer "+jwt)
		json.NewEncoder(w).Encode(map[string]interface{}{"jwt": jwt})
		return
	}

	w.WriteHeader(http.StatusForbidden)
}

func auth(name string, password string) bool {
	user := &models.User{Name: name}
	models.GetDB().Find(user)
	hash := lib.HashPassword(password, user.Salt)
	if hash == user.Password {
		return true
	}

	return false
}
