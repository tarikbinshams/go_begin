package main

import (
	"fmt"
	"log"
	"net/http"

	"hello/config"
	"hello/controllers"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()

	r.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	r.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")

	// Protected routes
	authRouter := r.PathPrefix("/users").Subrouter()
	authRouter.Use(controllers.AuthMiddleware)
	authRouter.HandleFunc("", controllers.CreateUser).Methods("POST")
	authRouter.HandleFunc("", controllers.GetUsers).Methods("GET")
	authRouter.HandleFunc("/{id}", controllers.GetUser).Methods("GET")
	authRouter.HandleFunc("/{id}", controllers.UpdateUser).Methods("PUT")
	authRouter.HandleFunc("/{id}", controllers.DeleteUser).Methods("DELETE")

	// Start Server
	fmt.Println("Server is running on port 4000...")
	log.Fatal(http.ListenAndServe(":4000", r))
}
