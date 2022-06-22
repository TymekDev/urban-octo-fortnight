package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router  *mux.Router
	storage Storage
}

var _ http.Handler = (*Server)(nil)

func NewServer(storage Storage) *Server {
	s := &Server{
		storage: storage,
	}
	r := mux.NewRouter()
	r.Methods("POST").Path("/user").HandlerFunc(s.userPOSTHandler)
	s.router = r
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type userPOSTPayload struct {
	Username string
}

func (s *Server) userPOSTHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var payload userPOSTPayload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := s.storage.NewUser(payload.Username); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
