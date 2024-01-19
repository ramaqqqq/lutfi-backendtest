package user

import (
	"os"

	"github.com/gorilla/mux"
)

func CustomerIdentityRouter(r *mux.Router) {

	middleUrl := os.Getenv("MIDDLE_URL")
	r.HandleFunc(middleUrl+"/user", userController.CreateUser).Methods("POST")
	r.HandleFunc(middleUrl+"/user", userController.GetList).Methods("GET")
	r.HandleFunc(middleUrl+"/user/{id}", userController.GetDetail).Methods("GET")
	r.HandleFunc(middleUrl+"/user/{id}", userController.UpdateUser).Methods("PUT")
	r.HandleFunc(middleUrl+"/user/{id}", userController.DeleteUser).Methods("DELETE")

	return
}
