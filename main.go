package main

import (
	"log"
	"net/http"

	"github.com/emj365/account/controllers"
	"github.com/emj365/account/models"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.PostUsers).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/auth", controllers.AuthUser).Methods("POST")
	log.Println("server is running on 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
	defer models.CloseDB()
}
