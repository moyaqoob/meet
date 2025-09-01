package auth

import(
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func Hander() http.Handler{
	handle:= mux.NewRouter()

	handle.HandleFunc("/auth/login",login);
	handle.HandleFunc("/auth/signup",signup)
	return handle
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is a signup func")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is a login func")
}
