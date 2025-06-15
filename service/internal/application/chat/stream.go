package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/gorilla/mux"
)

type ChatCompletionRequest struct {
	Model     string `json:"model"`
	Content   string `json:"content"`
	Reasoning int32  `json:"reasoning"`
}

func (s *Service) SendMessage(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	var body ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.log.Debug("failed to decode request body", "error", err)
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	chatID := mux.Vars(r)["id"]
	messages := make([]chat.Message, 0)
	rows, err := s.db.Query("SELECT role, content, reasoning FROM messages WHERE chat_id = ? AND user_id = ? ORDER BY created_at ASC", chatID, userID)
	if err != nil {
		s.log.Debug("failed to send query to database", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: Improve error handling
		return
	}
	defer rows.Close()

	for rows.Next() {
		var message chat.Message
		if err := rows.Scan(&message.Role, &message.Content, &message.Reasoning); err != nil {
			s.log.Debug("failed to scan message from database", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: Improve error handling
			return
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		s.log.Debug("failed to get chat history from database", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: Improve error handling
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
		Messages:            messages, // TODO: consider to split history and new message for better compatibility
	}

	stream, err := s.mr.StreamCompletion(req)
	if err != nil {
		s.log.Warn("failed to start a stream", "error", err)
		http.Error(w, "failed to start a stream", http.StatusInternalServerError)
		return
	}

	// Add stream to stream pool and return id
	// streamID := fmt.Sprintf("stream_%d", time.Now().UnixNano())
	streamID := fmt.Sprintf("stream_123")
	s.sp.Add(streamID, stream)

	s.log.Debug("stream was started sucessfully", "chat_id", chatID, "stream_id", streamID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"stream_id": streamID,
	}); err != nil {
		s.log.Error("failed to encode response", "error", err)
	}

}

func (s *Service) GetStream(w http.ResponseWriter, r *http.Request) {

	streamID := mux.Vars(r)["id"]

	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		s.log.Debug("streaming not supported", "stream_id", streamID)
		http.Error(w, "streaming not supported", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Subscribe to the stream
	sub, ok := s.sp.Subscribe(streamID)
	if !ok {
		s.log.Debug("stream not found", "stream_id", streamID)
		http.Error(w, "stream not found", http.StatusNotFound)
		return
	}

	// Flush all received chunks to the client
	for c := range sub {
		fmt.Printf("%v: receive chunk\n", time.Now().UnixMicro())
		fmt.Fprint(w, "event: message_delta\ndata: ")
		if err := json.NewEncoder(w).Encode(c); err != nil {
			panic("json encoding failed: " + err.Error())
		}
		fmt.Fprint(w, "\n")
		flusher.Flush()
	}
	fmt.Fprintf(w, "event: message_end\ndata: {\"done\":true}\n\n")
	flusher.Flush()

	s.log.Debug("streaming completed successfully", "stream_id", streamID)

}
