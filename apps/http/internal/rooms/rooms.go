package rooms

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Handler() http.Handler {
	handle := mux.NewRouter()
	handle.HandleFunc("rooms/create", createRoom)
	handle.HandleFunc("rooms/join", joinRoom)
	handle.HandleFunc("/rooms/leave", leaveRoom)

	return handle
}

func createRoom(w http.ResponseWriter, r *http.Request){
	
}

func joinRoom(w http.ResponseWriter, r *http.Request){
	
}

func leaveRoom(w http.ResponseWriter, r *http.Request){
	
}