package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/emj365/go-simple-account-service/controllers"
	"github.com/emj365/go-simple-account-service/libs"
	"github.com/emj365/go-simple-account-service/models"
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
	router.Use(jwtAuthentication) //attach JWT auth middleware
	router.HandleFunc("/users", controllers.PostUsers).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/me", controllers.GetMe).Methods("GET")
	router.HandleFunc("/auth", controllers.Auth).Methods("POST")
	router.HandleFunc("/jwt", controllers.JWT).Methods("GET")
	log.Println("server is running on 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
	defer models.CloseDB()
}

var jwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type requestInfo struct {
			Method string
			Path   string
		}

		notAuthRequests := []requestInfo{
			requestInfo{"POST", "/auth"},
			requestInfo{"POST", "/users"},
		} //List of endpoints that doesn't require auth

		requestPath := r.URL.Path //current request path
		requestMethod := r.Method

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuthRequests {
			if value.Method == requestMethod && value.Path == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		authorizationHeader := r.Header.Get("Authorization")
		token := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		jsonWebToken := libs.JsonWebToken{SecretKey: libs.GetSecretKey(),
			Token:  token,
			Claims: jwt.StandardClaims{}}
		err := jsonWebToken.Decode()
		if err != nil {
			log.Printf("Something went wrong: %s", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		sub, err := strconv.Atoi(jsonWebToken.Claims.Subject)
		if err != nil {
			log.Printf("Something went wrong: %s", err)
			libs.ResonponseServerError(w)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", uint(sub))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
