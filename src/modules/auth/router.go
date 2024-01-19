package auth

import (
	"os"

	"github.com/gorilla/mux"
)

func AuthRouter(r *mux.Router) {

	middleUrl := os.Getenv("MIDDLE_URL")
	r.HandleFunc(middleUrl+"/login", authController.Login).Methods("POST")
	r.HandleFunc(middleUrl+"/register", authController.Register).Methods("POST")

	return
}
