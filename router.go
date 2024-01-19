package main

import (
	"net/http"

	"folkatech-customerIdentity/src/modules/auth"
	user "folkatech-customerIdentity/src/modules/user"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	user.CustomerIdentityRouter(router)
	auth.AuthRouter(router)
}
