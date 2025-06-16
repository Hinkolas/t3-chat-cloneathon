package application

import (
	"net/http"
	"time"
)

func (app *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		app.Logger.Info("Request received",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"query", r.URL.RawQuery,
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		app.Logger.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration.String(),
		)
	})
}
