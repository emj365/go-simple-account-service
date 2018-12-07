package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emj365/account/lib"
	"github.com/emj365/account/models"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", postUsers).Methods("POST")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/auth", authUser).Methods("POST")
	log.Println("server is running on 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
	defer models.CloseDB()
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "getUsers")

	users := []models.User{}
	models.GetDB().Select("Name").Find(&users)
	lib.Resonponse(w, http.StatusOK, users)
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "postUsers")

	user := models.User{}
	if !getUserFromRequest(w, r, &user) {
		return
	}

	var hashedPassword, salt string
	ch := make(chan bool)
	go checkUserUnique(ch, user.Name, w)
	go genHashedPassword(ch, user.Password, &hashedPassword, &salt, w)

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
	models.GetDB().Create(&createdUser)
	createdUser.Password = "*******"
	lib.Resonponse(w, http.StatusCreated, createdUser)
}

func checkUserUnique(ch chan bool, name string, w http.ResponseWriter) {
	count := 0
	models.GetDB().Model(models.User{}).Where(models.User{Name: name}).Count(&count)
	if count > 0 {
		lib.Resonponse(w, http.StatusConflict, map[string]interface{}{"name": name})
		ch <- false
		return
	}

	ch <- true
}

func genHashedPassword(ch chan bool, password string, hashedPassword *string, salt *string, w http.ResponseWriter) {
	var err error
	*hashedPassword, *salt, err = lib.GenPasswordHash(password)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		lib.ResonponseServerError(w)
		ch <- false
		return
	}

	ch <- true
}

func authUser(w http.ResponseWriter, r *http.Request) {
	defer lib.TimeTrack(time.Now(), "auth")

	user := models.User{}
	if !getUserFromRequest(w, r, &user) {
		return
	}

	foundUser := models.User{Name: user.Name}
	count := 0
	models.GetDB().Where(foundUser).Select("Name, Password, Salt").Find(&foundUser).Count(&count)
	if count == 0 || !lib.Auth(foundUser, user.Name, user.Password) {
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

func getUserFromRequest(w http.ResponseWriter, r *http.Request, user *models.User) bool {
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil || user.Name == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}
