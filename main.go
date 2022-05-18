package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jose827corrza/go-websockets/handlers"
	"github.com/jose827corrza/go-websockets/server"
)

func main() {
	err := godotenv.Load(".env") //Load the .env file

	if err != nil {
		log.Fatal(".env file could not be loaded :c")
	}

	PORT := os.Getenv("PORT") //this last PORt is the name of the var given i the .env file
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DB := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JwtSecret:   JWT_SECRET,
		DatabaseUrl: DB,
	})

	if err != nil {
		log.Fatal(err)
	}
	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
}
