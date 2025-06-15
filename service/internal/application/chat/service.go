package chat

import (
	"database/sql"
	"log/slog"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm"
)

type Service struct {
	cfg *application.Config
	log *slog.Logger
	db  *sql.DB
	mr  *llm.ModelRouter
}

// NewService creates a new Chat service according to the provided config
func NewService(app *application.App) (*Service, error) {

	// Return initialized service
	return &Service{
		cfg: &app.Config,
		log: app.Logger,
		db:  app.Database,
		mr:  llm.NewModelRouter(),
	}, nil

}

// TODO: Delete later when implementing model config file
func (s *Service) AddModel(key string, model llm.Model) {
	// TODO: add some kind of error handling if model already exists
	s.mr.AddModel(key, model)
}
