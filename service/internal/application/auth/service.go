package auth

import (
	"database/sql"
	"log/slog"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application"
)

type Service struct {
	cfg *application.Config
	log *slog.Logger
	db  *sql.DB
}

// NewService creates a new Auth service according to the provided config
func NewService(app *application.App) (*Service, error) {

	// Return initialized service
	return &Service{
		cfg: &app.Config,
		log: app.Logger,
		db:  app.Database,
	}, nil

}
