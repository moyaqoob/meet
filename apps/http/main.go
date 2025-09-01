package main

import (
	"log"
	"meet/internal/auth"
	"meet/internal/room"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/auth/", auth.Handler())
	r.Handle("/rooms/", room.Handler())

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Go + Gorilla Mux!"))
	}).Methods("POST")

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
