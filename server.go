package gowell

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func defaultHealthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

type Server struct {
	router         *mux.Router
	healthzHandler func(http.ResponseWriter, *http.Request)
}

func NewServer() *Server {
	s := &Server{
		router:         mux.NewRouter(),
		healthzHandler: defaultHealthzHandler,
	}

	return s
}

func (s *Server) SetHealthzHandler(h func(http.ResponseWriter, *http.Request)) {
	s.healthzHandler = h
}

func (s *Server) SetRouter(r *mux.Router) {
	s.router = r
}

func (s *Server) Start() {
	s.router.HandleFunc("/healthz", s.healthzHandler)
	log.Fatal(http.ListenAndServe(":8080", s.router))
}
