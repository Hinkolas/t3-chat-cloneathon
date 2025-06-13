package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/gorilla/mux"
)

type ChatCompletionRequest struct {
	Model string `json:"model"`
}

func (s *Service) ChatCompletion(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body ChatCompletionRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Errorf("failed to decode request body: %w", err).Error(),
		}); err != nil {
			panic("json encoding failed: " + err.Error())
		}
		return
	}

	req := chat.Request{
		Model:               body.Model,
		Temperature:         0,
		MaxCompletionTokens: 1024,
		TopP:                1.0,
		Stream:              true,
		Thinking:            0,
		Stop:                nil,
		Messages: []chat.Message{
			{
				Role:    "user",
				Content: "Hello Claude!",
			},
		},
	}

	id := mux.Vars(r)["id"]

	completion, err := s.mr.ChatCompletion(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		}); err != nil {
			panic("json encoding failed: " + err.Error())
		}
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"chat_id":       id,
		"completion_id": completion,
	}); err != nil {
		panic("json encoding failed: " + err.Error())
	}

}
