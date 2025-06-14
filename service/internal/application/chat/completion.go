package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/gorilla/mux"
)

type ChatCompletionRequest struct {
	Model     string `json:"model"`
	Content   string `json:"content"`
	Reasoning int32  `json:"reasoning"`
}

func (s *Service) StreamMessage(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

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

	chatID := mux.Vars(r)["id"]
	messages := make([]chat.Message, 0)
	rows, err := s.db.Query("SELECT role, content, reasoning FROM messages WHERE chat_id = ? AND user_id = ? ORDER BY created_at ASC", chatID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var message chat.Message
		if err := rows.Scan(&message.Role, &message.Content, &message.Reasoning); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages = append(messages, chat.Message{
		Role:    "user",
		Content: body.Content,
	})

	req := chat.Request{
		Model:               body.Model,
		Temperature:         0,
		MaxCompletionTokens: 1024,
		TopP:                1.0,
		Stream:              true,
		Reasoning:           body.Reasoning,
		Stop:                nil,
		Messages:            messages,
	}

	stream, err := s.mr.ChatCompletion(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		}); err != nil {
			panic("json encoding failed: " + err.Error())
		}
		return
	}

	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "streaming not supported",
		}); err != nil {
			panic("json encoding failed: " + err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)

	// Subscribe to the stream
	sub := stream.Subscribe(1)

	for c := range sub {
		fmt.Fprint(w, "data: ")
		if err := json.NewEncoder(w).Encode(c); err != nil {
			panic("json encoding failed: " + err.Error())
		}
		fmt.Fprint(w, "\n")
		flusher.Flush()
	}

	stream.Close()

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()

}
