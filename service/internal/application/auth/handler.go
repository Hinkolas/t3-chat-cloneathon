package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Service) Handle(r *mux.Router) {

	router := r.PathPrefix("/v1/auth").Subrouter()

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("Login Request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Not Implemented"})
	})

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("Logout Request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Not Implemented"})
	})

	router.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("Profile Request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Not Implemented"})
	})

}
