package meetings

import (
	"encoding/json"
	"log"
	"meet/internal/models"
	"net/http"
	"time"

	middleware "meet/internal"

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

func Handler(handle *mux.Router) {
	var s *Server
	db := s.db;

	handle.Handle("/create", middleware.AuthMiddleware(createMeet(db))).Methods("POST")
	handle.Handle("/join", middleware.AuthMiddleware(joinMeet(db))).Methods("POST")
	handle.Handle("/leave", middleware.AuthMiddleware(leaveMeet(db))).Methods("POST")
	handle.Handle("/liveUsers", middleware.AuthMiddleware(fetchParticipants(db))).Methods("POST")
}

func createMeet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(uint)
		var meet models.Meeting
		json.NewDecoder(r.Body).Decode(&meet)

		random := generateId()

		newMeet := models.Meeting{
			HostID:    userID,
			Title:     meet.Title,
			MeetingID: random,
			CreatedAt: time.Now(),
		}

		err := db.Create(&newMeet)
		if err != nil {
			log.Fatal("error creating a meeting")
			return
		}

		response := map[string]string{
			"title":     newMeet.Title,
			"meetingId": newMeet.MeetingID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}

}

func joinMeet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userId from context (set by AuthMiddleware)
		userId := r.Context().Value("userId").(uint)

		var meetUser models.MeetingParticipant
		if err := json.NewDecoder(r.Body).Decode(&meetUser); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		var meet models.Meeting
		if err := db.Where("meeting_id = ?", meetUser.MeetingID).First(&meet).Error; err != nil {
			http.Error(w, "meeting not found", http.StatusNotFound)
			return
		}

		meetUser.UserID = userId
		if err := db.Create(&meetUser).Error; err != nil {
			http.Error(w, "could not join", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("joined meeting successfully"))
	}
}

func leaveMeet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userId").(uint)

		var participation models.MeetingParticipant
		err := db.Where("user_id = ? AND active = ?", userID, true).First(&participation).Error
		if err != nil {
			http.Error(w, "no active meeting found", http.StatusNotFound)
			return
		}

		meetingID := participation.MeetingID

		removeUser := db.Where("meetingID = ?", meetingID).Delete(&participation)
		if removeUser.Error != nil {
			http.Error(w, "remove the user", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("leave the meeting"))

	}
}

func fetchParticipants(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userId").(uint)


		// 1. Find active meeting(s) for this user
		var participation models.MeetingParticipant
		err := db.Where("user_id = ? AND is_active = ?", userID, true).First(&participation).Error
		if err != nil {
			http.Error(w, "no active meeting found", http.StatusNotFound)
			return
		}

		meetingID := participation.MeetingID

		// 2. Fetch all active participants in that meeting
		var participants []models.MeetingParticipant
		err = db.Where("meeting_id = ? AND is_active = ?", meetingID, true).Find(&participants).Error
		if err != nil {
			http.Error(w, "could not fetch participants", http.StatusInternalServerError)
			return
		}

		// 3. Return as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(participants)
	}
}
