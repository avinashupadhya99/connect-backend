package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func initializeRouter() {
	r := mux.NewRouter()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/api/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/api/users", CreateUser).Methods("POST")
	r.HandleFunc("/api/users/login", LoginUser).Methods("POST")
	r.HandleFunc("/api/slot/book", BookSlot).Methods("POST")
	r.HandleFunc("/api/slot", GetSlot).Methods("GET")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(credentials, methods, origins)(r)))
}

func main() {
	InitialMigration()
	initializeRouter()
}
