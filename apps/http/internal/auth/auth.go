package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"meet/internal/models"
	"net/http"
	"os"
	"time"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

type Server struct {
	DB *gorm.DB
}

func Handler(handle *mux.Router) {
	err := godotenv.Load("./.env")

	fmt.Println(err, "env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database := os.Getenv("DB_URL")
	fmt.Println(database, "database url")
	db, err := gorm.Open(postgres.Open(database), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	fmt.Println("connected database successfully")

	handle.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Go + Gorilla Mux!"))
	}).Methods("POST")
	server := &Server{DB: db}
	handle.HandleFunc("/login", server.login).Methods("POST")
	handle.HandleFunc("/signup", server.signup).Methods("POST")
	handle.HandleFunc("/test", test).Methods("POST")
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("writing form the auth.go"))
	fmt.Println("writing fromt he auth.go")
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I reached here")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	db := s.DB

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		fmt.Println("password could not hash")
	}
	fmt.Println("user", &user, "password", hashPassword)

	userExists := db.Where("email = ?", user.Email).First(&user)
	if userExists.Error != nil {
		fmt.Println("user exists")
	}

	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashPassword),
	}

	if err := s.DB.Create(&newUser).Error; err != nil {
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created successfully"))
	json.NewEncoder(w).Encode(newUser)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I reached here too")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	jwtKey := []byte(os.Getenv("SECRET_KEY"))
	var existingUser models.User
	result := s.DB.Where("email = ?", user.Email).First(&existingUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		// other DB errors
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	fmt.Println("here is the route")
	expirationTime := time.Now().Add(24 * time.Hour) // 24h token
	claims := &Claims{
		UserID: existingUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err!=nil{
		fmt.Println("Invalid password")
	}else{
		fmt.Println("Correct Password")
	}

	fmt.Println("expiration and bcrypt pass",claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token",token)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	fmt.Println("token created just returning the token")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
