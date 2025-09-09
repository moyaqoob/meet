package meetings

import (
	"encoding/json"
	"log"
	"meet/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func generateId() string {
	str := uuid.New().String()
	return str[:10]
}
 
type Server struct {
	db *gorm.DB
}

func Handler(handle *mux.Router, gormDB *gorm.DB) {
	var s *Server
	db :=s.db
	handle.HandleFunc("/create", createMeet(db)).Methods("POST")
	handle.HandleFunc("/join", joinMeet(db)).Methods("POST")
	handle.HandleFunc("/leave", leaveMeet(db)).Methods("POST")
}

func createMeet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(uint)
		var meet models.Meeting
		json.NewDecoder(r.Body).Decode(&meet)

		random := generateId()

		newMeet:= models.Meeting{
			HostID: userID,
			Title:meet.Title ,
			MeetingID: random,
			CreatedAt: time.Now(),
		}

		err := db.Create(&newMeet)
		if err!=nil{
			log.Fatal("error creating a meeting")
			return
		}

		response := map[string]string{
			"title":newMeet.Title,
			"meetingId":newMeet.MeetingID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}

}

func joinMeet(db *gorm.DB) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is join room from rooms.go "))
	}

}
func leaveMeet(db *gorm.DB) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is join room from rooms.go "))
	}

}