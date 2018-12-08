package main

import (
	"log"
	"net/http"

	"github.com/emj365/account/controllers"
	"github.com/emj365/account/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// https://github.com/joho/godotenv#usage
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.PostUsers).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/auth", controllers.AuthUser).Methods("POST")
	router.HandleFunc("/jwt", controllers.JWT).Methods("GET")
	log.Println("server is running on 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
	defer models.CloseDB()
}
