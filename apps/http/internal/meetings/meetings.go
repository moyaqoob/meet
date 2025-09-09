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

func Handler(handle *mux.Router) {
	handle.HandleFunc("/create", createMeet).Methods("POST")
	handle.HandleFunc("/join", joinMeet).Methods("POST")
	handle.HandleFunc("/leave", leaveMeet).Methods("POST")
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

func joinMeet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is join room from rooms.go "))
}

func leaveMeet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is leave room from rooms.go "))
}
