package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
	"github.com/gorilla/mux"
)

type ChatCompletionRequest struct {
	// Message
	Content     string   `json:"content"`
	Attachments []string `json:"attachments,omitempty"`
	// Options
	Model           string `json:"model"`
	ReasoningEffort int32  `json:"reasoning_effort"`
}

func (s *Service) AddMessage(w http.ResponseWriter, r *http.Request) {

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
		MaxCompletionTokens: 8192,
		TopP:                1.0,
		Stream:              true,
		ReasoningEffort:     body.ReasoningEffort,
		Stop:                nil,
		Messages:            messages, // TODO: Consider to split history and new message for better compatibility
	}

	now := time.Now()
	message := Message{
		ID:        fmt.Sprintf("msg_%d", now.UnixNano()),
		ChatID:    chatID,
		UserID:    userID,
		Role:      "user",
		Status:    "done",
		Model:     req.Model,
		Content:   body.Content,
		CreatedAt: now.UnixMilli(),
		UpdatedAt: now.UnixMilli(),
	}

	_, err = s.db.Exec("INSERT INTO messages (id, chat_id, user_id, stream_id, role, status, model, content, reasoning, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		message.ID, message.ChatID, message.UserID, message.StreamID,
		message.Role, message.Status, message.Model,
		message.Content, message.Reasoning,
		message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		s.log.Warn("failed to insert message into database", "error", err)
		fmt.Println(err)
		http.Error(w, "failed to insert message into database", http.StatusInternalServerError)
		return
	}

	compl, err := s.mr.StreamCompletion(req)
	if err != nil {
		s.log.Warn("failed to start a stream", "error", err)
		http.Error(w, "failed to start a stream", http.StatusInternalServerError)
		return
	}

	now = time.Now()
	streamID := fmt.Sprintf("stream_%d", now.UnixNano()) // TODO: Consider replacing with uuid
	message = Message{
		ID:        fmt.Sprintf("msg_%d", now.UnixNano()),
		ChatID:    chatID,
		UserID:    userID,
		StreamID:  streamID,
		Role:      "assistant",
		Status:    "streaming",
		Model:     req.Model,
		CreatedAt: now.UnixMilli(),
		UpdatedAt: now.UnixMilli(),
	}

	_, err = s.db.Exec("INSERT INTO messages (id, chat_id, user_id, stream_id, role, status, model, content, reasoning, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		message.ID, message.ChatID, message.UserID, message.StreamID,
		message.Role, message.Status, message.Model,
		message.Content, message.Reasoning,
		message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		s.log.Warn("failed to insert message into database", "error", err)
		fmt.Println(err)
		http.Error(w, "failed to insert message into database", http.StatusInternalServerError)
		return
	}

	// Add stream to stream pool and return id
	s.sp.Add(streamID, compl)
	compl.OnClose(func(chunks []stream.Chunk) error {

		var contentBuilder, reasoningBuilder strings.Builder

		for _, chunk := range chunks {
			contentBuilder.WriteString(chunk.Content)
			reasoningBuilder.WriteString(chunk.Reasoning)
		}

		_, err = s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
			contentBuilder.String(), reasoningBuilder.String(), "done", time.Now().UnixMilli(), message.ID,
		)

		return err

	})

	s.log.Debug("stream was started sucessfully", "chat_id", chatID, "stream_id", streamID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"stream_id": streamID,
	}); err != nil {
		s.log.Error("failed to encode response", "error", err)
	}

}

func (s *Service) SendMessage(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	var body ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.log.Debug("failed to decode request body", "error", err)
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	chatID := fmt.Sprintf("chat_%d", time.Now().UnixNano())
	newChat := Chat{
		ID:            chatID,
		UserID:        userID,
		Title:         fmt.Sprintf("New Chat %d", time.Now().Unix()), // TODO: Generate title based on the first message using a lightweight model
		Model:         body.Model,
		IsPinned:      false,
		Status:        "streaming",
		CreatedAt:     now.UnixMilli(),
		UpdatedAt:     now.UnixMilli(),
		LastMessageAt: 0,
	}

	_, err := s.db.Exec("INSERT INTO chats (id, user_id, title, model, is_pinned, status, created_at, updated_at, last_message_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newChat.ID, newChat.UserID,
		newChat.Title, newChat.Model,
		newChat.IsPinned, newChat.Status,
		newChat.CreatedAt, newChat.UpdatedAt,
		newChat.LastMessageAt,
	)
	if err != nil {
		s.log.Warn("failed to insert chat into database", "error", err)
		fmt.Println(err)
		http.Error(w, "failed to insert chat into database", http.StatusInternalServerError)
		return
	}

	messages := []chat.Message{
		{
			Role:    "user",
			Content: body.Content,
		},
	}

	req := chat.Request{
		Model:               body.Model,
		Temperature:         0,
		MaxCompletionTokens: 1024,
		TopP:                1.0,
		Stream:              true,
		ReasoningEffort:     body.ReasoningEffort,
		Stop:                nil,
		Messages:            messages, // TODO: Consider to split history and new message for better compatibility
	}

	now = time.Now()
	message := Message{
		ID:        fmt.Sprintf("msg_%d", now.UnixNano()),
		ChatID:    newChat.ID,
		UserID:    userID,
		Role:      "user",
		Status:    "done",
		Model:     req.Model,
		Content:   body.Content,
		CreatedAt: now.UnixMilli(),
		UpdatedAt: now.UnixMilli(),
	}

	_, err = s.db.Exec("INSERT INTO messages (id, chat_id, user_id, stream_id, role, status, model, content, reasoning, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		message.ID, message.ChatID, message.UserID, message.StreamID,
		message.Role, message.Status, message.Model,
		message.Content, message.Reasoning,
		message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		s.log.Warn("failed to insert message into database", "error", err)
		fmt.Println(err)
		http.Error(w, "failed to insert message into database", http.StatusInternalServerError)
		return
	}

	compl, err := s.mr.StreamCompletion(req)
	if err != nil {
		s.log.Warn("failed to start a stream", "error", err)
		http.Error(w, "failed to start a stream", http.StatusInternalServerError)
		return
	}

	now = time.Now()
	streamID := fmt.Sprintf("stream_%d", now.UnixNano()) // TODO: Consider replacing with uuid
	message = Message{
		ID:        fmt.Sprintf("msg_%d", now.UnixNano()),
		ChatID:    chatID,
		UserID:    userID,
		StreamID:  streamID,
		Role:      "assistant",
		Status:    "streaming",
		Model:     req.Model,
		CreatedAt: now.UnixMilli(),
		UpdatedAt: now.UnixMilli(),
	}

	_, err = s.db.Exec("INSERT INTO messages (id, chat_id, user_id, stream_id, role, status, model, content, reasoning, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		message.ID, message.ChatID, message.UserID, message.StreamID,
		message.Role, message.Status, message.Model,
		message.Content, message.Reasoning,
		message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		s.log.Warn("failed to insert message into database", "error", err)
		fmt.Println(err)
		http.Error(w, "failed to insert message into database", http.StatusInternalServerError)
		return
	}

	// Add stream to stream pool and return id
	s.sp.Add(streamID, compl)
	compl.OnClose(func(chunks []stream.Chunk) error {

		var contentBuilder, reasoningBuilder strings.Builder

		for _, chunk := range chunks {
			contentBuilder.WriteString(chunk.Content)
			reasoningBuilder.WriteString(chunk.Reasoning)
		}

		_, err = s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
			contentBuilder.String(), reasoningBuilder.String(), "done", time.Now().UnixMilli(), message.ID,
		)

		return err

	})

	s.log.Debug("stream was started sucessfully", "chat_id", newChat.ID, "stream_id", streamID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"chat_id":   newChat.ID,
		"stream_id": streamID,
	}); err != nil {
		s.log.Error("failed to encode response", "error", err)
	}

}
