package chat

import (
	"github.com/gorilla/mux"
)

func (s *Service) Handle(r *mux.Router) {

	router := r.PathPrefix("/v1/chats").Subrouter()

	router.HandleFunc("/", s.ListChats).Methods("GET")
	router.HandleFunc("/{id}/", s.GetChat).Methods("GET")
	router.HandleFunc("/{id}/", s.DeleteChat).Methods("DELETE")
	router.HandleFunc("/{id}/", s.EditChat).Methods("PATCH")
	router.HandleFunc("/{id}/", s.ChatCompletion).Methods("POST")

	// TODO: Maybe move somewhere else
	router = r.PathPrefix("/v1/models").Subrouter()

	router.HandleFunc("/", s.ListModels).Methods("GET")

}
