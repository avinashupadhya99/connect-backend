package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func initializeRouter() {
	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/api/users", GetUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/api/users/{id}/meetings", GetUserMeetings).Methods("GET")
	r.HandleFunc("/api/users", CreateUser).Methods("POST")
	r.HandleFunc("/api/users/login", LoginUser).Methods("POST")
	r.HandleFunc("/api/slot/book", BookSlot).Methods("POST")
	r.HandleFunc("/api/slot/unbook", UnBookSlot).Methods("DELETE")
	r.HandleFunc("/api/slot", GetSlot).Methods("GET")
	r.HandleFunc("/api/slot", DeleteSlot).Methods("DELETE")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(r)))
}

func InitializeCron() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Days().Do(createMeetings)
	s.StartAsync()
}

func main() {
	InitialMigration()
	InitializeCron()
	initializeRouter()
}
