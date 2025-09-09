package main

import (
	"log"
	"meet/internal/auth"
	"meet/internal/meetings"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
   
	// Mount /auth routes
	authRouter := r.PathPrefix("/auth").Subrouter()
	auth.Handler(authRouter)
	meetingRouter :=  r.PathPrefix("/meeting").Subrouter()
	meetings.Handler(meetingRouter)

	// Example root route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Go + Gorilla Mux!"))
	}).Methods("GET")

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}