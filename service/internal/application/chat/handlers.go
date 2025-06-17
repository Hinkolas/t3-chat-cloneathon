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
	router.HandleFunc("/", s.SendMessage).Methods("POST")
	router.HandleFunc("/{id}/", s.GetChat).Methods("GET")
	router.HandleFunc("/{id}/", s.DeleteChat).Methods("DELETE")
	router.HandleFunc("/{id}/", s.EditChat).Methods("PATCH")
	router.HandleFunc("/{id}/", s.AddMessage).Methods("POST")

	router = r.PathPrefix("/v1/attachments").Subrouter()
	router.HandleFunc("/", s.ListAttachments).Methods("GET")
	router.HandleFunc("/", s.UploadAttachment).Methods("POST")
	router.HandleFunc("/{id}/", s.GetAttachment).Methods("GET")
	router.HandleFunc("/{id}/", s.DeleteAttachment).Methods("DELETE")

	router = r.PathPrefix("/v1/streams").Subrouter()
	router.HandleFunc("/{id}/", s.OpenStream).Methods("GET")
	router.HandleFunc("/{id}/", s.CancelStream).Methods("DELETE")

	router = r.PathPrefix("/v1/profile").Subrouter()
	router.HandleFunc("/", s.GetUserProfile).Methods("GET")
	router.HandleFunc("/", s.UpsertUserProfile).Methods("PATCH")

	router = r.PathPrefix("/v1/share").Subrouter()
	router.HandleFunc("/{id}/", s.GetSharedChat).Methods("GET")

}
