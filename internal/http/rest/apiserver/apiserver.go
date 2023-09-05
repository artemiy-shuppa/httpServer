package apiserver

import (
	"fmt"
	"log/slog"
	"myHttpServer/internal/http/rest/handlers"
	"myHttpServer/internal/postgres"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type ApiServer struct {
	router *mux.Router
	addr   string

	log            *slog.Logger
	db             *postgres.Database
	messageHandler *handlers.MessageHandler
}

func New(log *slog.Logger) (*ApiServer, error) {

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	dataSourceName := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=app",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
	)
	db, err := postgres.Open(dataSourceName)
	if err != nil {
		log.Error("Can't connect to database")
		return nil, err
	}
	messageStore := db.MessageStorage()
	if err != nil {
		log.Error("Can't find required tables")
		return nil, err
	}
	return &ApiServer{
		log:            log,
		router:         mux.NewRouter(),
		addr:           os.Getenv("ADDR"),
		messageHandler: handlers.NewMessageHandler(messageStore, log),
	}, nil

}

func (s *ApiServer) Start() error {
	s.configureRouter()

	s.log.Info("Server started")
	return http.ListenAndServe(s.addr, s.router)
}

func (s *ApiServer) configureRouter() {
	getRouter := s.router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/messages", s.messageHandler.Messages())
	getRouter.HandleFunc("/messages/{id:[0-9]+}", s.messageHandler.Message())

	postRouter := s.router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/messages", s.messageHandler.CreateMessage())

	putRouter := s.router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/messages/{id:[0-9]+}", s.messageHandler.UpdateMessage())

	deleteRouter := s.router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/messages/{id:[0-9]+}", s.messageHandler.DeleteMessage())
}
