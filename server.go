package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router  *mux.Router
	storage Storage
}

var _ http.Handler = (*Server)(nil)

func NewServer(storage Storage) *Server {
	return &Server{
		router: mux.NewRouter
		storage: storage,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
