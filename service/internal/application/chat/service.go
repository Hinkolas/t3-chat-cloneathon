package chat

import (
	"database/sql"
	"log/slog"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
)

type Service struct {
	cfg *application.Config
	log *slog.Logger
	db  *sql.DB
	mr  *llm.ModelRouter
	sp  *stream.StreamPool
}

// NewService creates a new Chat service according to the provided config
func NewService(app *application.App) (*Service, error) {

	// Return initialized service
	return &Service{
		cfg: &app.Config,
		log: app.Logger,
		db:  app.Database,
		mr:  llm.NewModelRouter(),
		sp:  stream.NewStreamPool(),
	}, nil

}
