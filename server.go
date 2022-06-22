package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	game   Game
}

var _ http.Handler = (*Server)(nil)

func NewServer(game Game) *Server {
	s := &Server{
		game: game,
	}
	r := mux.NewRouter()
	r.Methods("POST").Path("/user").HandlerFunc(s.userPOSTHandler)
	r.Methods("GET").Path("/dashboard").HandlerFunc(s.dashboardGETHandler)
	r.Methods("POST").Path("/upgrade").HandlerFunc(s.upgradePOSTHandler)
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
	if err := s.game.NewUser(payload.Username); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// Consciously not making the wrong abstraction for possible extendability.
// Reference: https://szymanskir.github.io/post/2022-04-30-the-problem-with-dry/
type dashboardGETPayload struct {
	Username string
}

type dashboardGETResponse struct {
	Resources Resources
	Factories Factories
}

func (s *Server) dashboardGETHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var payload dashboardGETPayload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	userData, err := s.game.GetUserData(payload.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	response := dashboardGETResponse{
		Resources: userData.Resources(),
		Factories: userData.Factories(),
	}
	bytes, err = json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

type upgradePOSTPayload struct {
	Username string
	Factory  string
}

func (s *Server) upgradePOSTHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var payload upgradePOSTPayload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	factoryType, err := FactoryTypeFromString(payload.Factory)
	if err := json.Unmarshal(bytes, &payload); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := s.game.UpgradeFactory(payload.Username, factoryType); err != nil {
		log.Println(err)
		// TODO: this could be improved because it is either missing resources or nonexistent user
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}
