package auth

import (
	"context"
	"net/http"
	"strings"
)

func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			s.log.Debug("No authorization header provided")
			next.ServeHTTP(w, r)
			return
		}

		// Split by space and check the format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			s.log.Debug("Invalid authorization header format", "auth_header", authHeader)
			next.ServeHTTP(w, r)
			return
		}

		sessionToken := parts[1]
		if len(sessionToken) == 0 {
			s.log.Debug("No session token provided", "auth_header", authHeader)
			next.ServeHTTP(w, r)
			return
		}

		session, err := s.Authorize(r.Context(), sessionToken)
		if err != nil {
			s.log.Debug("Failed to authorize session", "session_token", sessionToken, "error", err)
			next.ServeHTTP(w, r)
			return
		}

		// Add UserID and SessionToken to the request context
		s.log.Info("Authorized session", "user_id", session.UserID)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", session.UserID)
		ctx = context.WithValue(ctx, "session_id", session.ID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
