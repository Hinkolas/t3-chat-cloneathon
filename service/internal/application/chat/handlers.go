package chat

import (
	"github.com/gorilla/mux"
)

func (s *Service) Handle(r *mux.Router) {

	router := r.PathPrefix("/v1/chats").Subrouter()

	router.HandleFunc("/", s.GetChats).Methods("GET")
	router.HandleFunc("/{id}", s.GetChat).Methods("GET")
	router.HandleFunc("/{id}", s.GetChat).Methods("POST")

	// TODO: Maybe move somewhere else
	router = r.PathPrefix("/v1/models").Subrouter()

	router.HandleFunc("/", s.GetModels).Methods("GET")

}
