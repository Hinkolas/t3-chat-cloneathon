package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Service) Handle(r *mux.Router) {

	r.Use(s.AuthMiddleware)

	router := r.PathPrefix("/v1/auth").Subrouter()

	router.HandleFunc("/login/", s.Login)
	router.HandleFunc("/logout/", s.Logout)
	router.HandleFunc("/session/", s.GetCurrentSession)

}

// LoginRequest represents a request to log in to the system.
type LoginRequest struct {
	// Username and password for authentication.
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if ok {
		s.log.Debug("User is already authenticated", "user_id", userID)
		http.Error(w, "already_authenticated", http.StatusConflict)
		return
	}

	// Decode the request body into the required request format
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.log.Debug("Invalid login request payload", "error", err)
		http.Error(w, "invalid_request_payload", http.StatusBadRequest)
		return
	}

	user, err := s.Authenticate(r.Context(), request.Username, request.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			s.log.Debug("Invalid credentials", "username", request.Username, "error", err)
			http.Error(w, "invalid_credentials", http.StatusUnauthorized)
		} else {
			s.log.Warn("Failed to authenticate user", "username", request.Username, "error", err)
			http.Error(w, "failed_authenticate", http.StatusInternalServerError)
		}
		return
	}

	session, err := s.CreateSession(r.Context(), user)
	if err != nil {
		s.log.Warn("Failed to create session", "user_id", user.ID)
		http.Error(w, "failed_create_session", http.StatusInternalServerError)
		return
	}

	// Instruct the client to store the session token in a cookie
	cookie := &http.Cookie{
		Name:        "session_token",
		Value:       session.Token,
		Path:        "/",
		HttpOnly:    true,
		Expires:     time.Unix(session.IssuedAt+session.TimeToLive, 0), // Set to the Unix epoch (very past date)
		Partitioned: true,
		SameSite:    http.SameSiteNoneMode, // TODO: Replace with somehting more secure
		// TODO: add secure in prod
	}
	http.SetCookie(w, cookie)

	// Send success status and session details to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"session": session,
	})

}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	sessionID, ok := r.Context().Value("session_id").(uuid.UUID)
	if !ok {
		s.log.Debug("Unautorized request to logout handler")
		http.Error(w, "not_authorized", http.StatusBadRequest)
		return
	}

	// Insert new session into database
	err := s.DeleteSession(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "failed_delete_session", http.StatusInternalServerError)
		return
	}

	// Instruct the client to store the session token in a cookie
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0), // Set to the Unix epoch (very past date)
	}
	http.SetCookie(w, cookie)

	// Send success status to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{})

}

func (s *Service) GetCurrentSession(w http.ResponseWriter, r *http.Request) {

	// Get sessionID from auth middleware, ok if authenticated
	sessionID, ok := r.Context().Value("session_id").(uuid.UUID)
	if !ok {
		s.log.Debug("Unauthorized request to get current session handler")
		http.Error(w, "not_authorized", http.StatusUnauthorized)
		return
	}

	session, err := s.GetSession(r.Context(), sessionID)
	if err != nil {
		s.log.Warn("Failed to get session", "session_id", sessionID, "error", err)
		http.Error(w, "failed_get_session", http.StatusInternalServerError)
		return
	}

	// Send success status and session details to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"session": session,
	})

}
