package chat

import (
	"encoding/json"
	"net/http"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm"
)

// TODO: Delete later when implementing model config file
func (s *Service) AddModel(key string, model llm.Model) {
	// TODO: add some kind of error handling if model already exists
	s.mr.AddModel(key, model)
}

func (s *Service) ListModels(w http.ResponseWriter, r *http.Request) {

	models := s.mr.ListModels()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
