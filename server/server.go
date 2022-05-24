package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jose827corrza/go-websockets/database"
	"github.com/jose827corrza/go-websockets/repository"
	"github.com/jose827corrza/go-websockets/websocket"
	cors "github.com/rs/cors"
)

type Config struct {
	Port        string
	JwtSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port is required")
	}
	if config.JwtSecret == "" {
		return nil, errors.New("JWT is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("DB URL is required")
	}
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	b.router = mux.NewRouter()
	binder(b, b.router)
	handler := cors.Default().Handler(b.router) //Importante que este despues de la funcion binder -.-
	log.Println("Server running on port:", b.Config().Port)
	if err := http.ListenAndServe(b.Config().Port, handler); err != nil { //El segundo parametro "b.router"
		//Se cambia por el handler, esto para quitar el error CORS
		log.Fatal("ListenAndServe:", err)
	}
}
