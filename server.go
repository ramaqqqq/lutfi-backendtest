package main

import (
	"fmt"
	"folkatech-customerIdentity/src/config"
	"folkatech-customerIdentity/src/middleware"
	"folkatech-customerIdentity/src/modules/auth"
	"folkatech-customerIdentity/src/modules/user"
	"folkatech-customerIdentity/src/pkg/db"
	"folkatech-customerIdentity/src/pkg/helpers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func InitApp() {
	err := godotenv.Load()
	if err != nil {
		helpers.Logger("error", "Error getting env")
	}

	cfg := config.NewConfig()
	db := db.NewDbConnection(cfg)

	userModule := user.New(db.Mongo, *db.Redis)
	userModule.InitModule()
	auth.New().InitModule()

	r := mux.NewRouter()
	r.Use(middleware.JwtAuth)
	SetupRoutes(r)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "content-type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	p := os.Getenv("PORT")
	h := c.Handler(r)
	s := new(http.Server)
	s.Handler = h
	s.Addr = ":" + p

	fmt.Printf("Server listening on port %s\n", p)
	s.ListenAndServe()
}
