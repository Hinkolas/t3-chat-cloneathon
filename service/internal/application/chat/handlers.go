package chat

import (
	"github.com/gorilla/mux"
)

func (s *Service) Handle(r *mux.Router) {

	var router *mux.Router

	router = r.PathPrefix("/v1/models").Subrouter()
	router.HandleFunc("/", s.ListModels).Methods("GET")

	router = r.PathPrefix("/v1/chats").Subrouter()
	router.HandleFunc("/", s.ListChats).Methods("GET")
	router.HandleFunc("/{id}/", s.GetChat).Methods("GET")
	router.HandleFunc("/{id}/", s.DeleteChat).Methods("DELETE")
	router.HandleFunc("/{id}/", s.EditChat).Methods("PATCH")
	router.HandleFunc("/{id}/", s.StreamMessage).Methods("POST")

	router = r.PathPrefix("/v1/attachments").Subrouter()
	router.HandleFunc("/", s.ListAttachments).Methods("GET")
	// router.HandleFunc("/{id}/", s.GetAttachment).Methods("GET")
	router.HandleFunc("/{id}/", s.DeleteAttachment).Methods("DELETE")
	// router.HandleFunc("/{id}/", s.UploadAttachment).Methods("POST")

}
