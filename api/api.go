package api

import (
	"database/sql"
	"go-note/service/auth"
	"go-note/service/note"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := auth.NewStore(s.db)
	userHandler := auth.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	noteStore := note.NewStore(s.db)
	noteHandler := note.NewHandler(noteStore)
	noteHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
