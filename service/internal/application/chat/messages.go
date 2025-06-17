package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

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

	profile, err := s.getUserProfile(userID)
	if err != nil {
		s.log.Warn("failed to get user profile", "error", err)
	}

	now := time.Now()
	message := Message{
		ID:        fmt.Sprintf("msg_%d", now.UnixNano()),
		ChatID:    chatID,
		UserID:    userID,
		Role:      "user",
		Status:    "done",
		Model:     body.Model,
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

	attachments, err := s.loadAndUpdateAttachments(body.Attachments, message.ID, userID)
	if err != nil {
		s.log.Warn("failed to attach all attachments to message", "message_id", message.ID, "error", err)
		// Note: Consider whether this should be a fatal error or just logged
	}

	req := chat.Request{
		Model:               body.Model,
		Temperature:         0.7,
		MaxCompletionTokens: 8192,
		TopP:                1.0,
		Stream:              true,
		ReasoningEffort:     body.ReasoningEffort,
		Stop:                nil,
		Messages:            append(messages, chat.Message{Role: "user", Content: body.Content, Attachments: attachments}), // TODO: Consider to split history and new message for better compatibility
		System:              profile.SystemPrompt(),
	}

	compl, err := s.mr.StreamCompletion(req, profile.Options())
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
	compl.OnClose(func(chunk stream.Chunk, serr error) {

		if serr != nil {
			s.log.Error("stream failed", "stream_id", streamID, "error", serr)
			return
		}

		_, err := s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
			chunk.Content, chunk.Reasoning, "done", time.Now().UnixMilli(), message.ID,
		)

		if err != nil {
			s.log.Error("storing stream content failed", "stream_id", streamID, "error", err)
			return
		}

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

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

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
		SharedAt:      0,
	}

	_, err := s.db.Exec("INSERT INTO chats (id, user_id, title, model, is_pinned, status, created_at, updated_at, last_message_at, shared_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newChat.ID, newChat.UserID,
		newChat.Title, newChat.Model,
		newChat.IsPinned, newChat.Status,
		newChat.CreatedAt, newChat.UpdatedAt,
		newChat.LastMessageAt, newChat.SharedAt,
	)
	if err != nil {
		s.log.Warn("failed to insert chat into database", "error", err)
		fmt.Println(err)
		http.Error(w, "failed to insert chat into database", http.StatusInternalServerError)
		return
	}

	profile, err := s.getUserProfile(userID)
	if err != nil {
		s.log.Warn("failed to get user profile", "error", err)
	}

	now = time.Now()
	message := Message{
		ID:        fmt.Sprintf("msg_%d", now.UnixNano()),
		ChatID:    newChat.ID,
		UserID:    userID,
		Role:      "user",
		Status:    "done",
		Model:     body.Model,
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

	attachments, err := s.loadAndUpdateAttachments(body.Attachments, message.ID, userID)
	if err != nil {
		s.log.Warn("failed to attach all attachments to message", "message_id", message.ID, "error", err)
		// Note: Consider whether this should be a fatal error or just logged
	}

	req := chat.Request{
		Model:               body.Model,
		Temperature:         0.7,
		MaxCompletionTokens: 8192,
		TopP:                1.0,
		Stream:              true,
		ReasoningEffort:     body.ReasoningEffort,
		Stop:                nil,
		Messages:            []chat.Message{{Role: "user", Content: body.Content, Attachments: attachments}}, // TODO: Consider to split history and new message for better compatibility
		System:              profile.SystemPrompt(),
	}

	compl, err := s.mr.StreamCompletion(req, profile.Options())
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
	compl.OnClose(func(chunk stream.Chunk, serr error) {

		if serr != nil {
			s.log.Error("stream failed", "stream_id", streamID, "error", serr)
			return
		}

		_, err := s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
			chunk.Content, chunk.Reasoning, "done", time.Now().UnixMilli(), message.ID,
		)

		if err != nil {
			s.log.Error("storing stream content failed", "stream_id", streamID, "error", err)
			return
		}

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

// LoadAndUpdateAttachments updates message_id on the given attachments and returns their mime_type+data.
func (s *Service) loadAndUpdateAttachments(attachmentIDs []string, messageID, userID string) ([]chat.Attachment, error) {

	if len(attachmentIDs) == 0 {
		return []chat.Attachment{}, nil
	}

	// 1) Build the IN-clause placeholders and args:
	placeholders := make([]string, len(attachmentIDs))
	args := make([]any, 0, len(attachmentIDs)+2)

	// first arg is the new message_id
	args = append(args, messageID)
	// next come the attachment IDs
	for i, id := range attachmentIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}
	// last arg is the user_id
	args = append(args, userID)

	// Update the attachments to link to the message and return src+mime_type:
	query := fmt.Sprintf("UPDATE attachments SET message_id = ? WHERE id IN (%s) AND user_id = ? AND message_id = '' RETURNING id, type", strings.Join(placeholders, ","))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("update attachments with returning: %w", err)
	}
	defer rows.Close()

	// Scan out the returned rows into chat.Attachment
	var attachments []Attachment
	for rows.Next() {
		var attachment Attachment
		if err := rows.Scan(&attachment.ID, &attachment.Type); err != nil {
			return nil, fmt.Errorf("scan returned attachment: %w", err)
		}
		attachments = append(attachments, attachment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating returned attachments: %w", err)
	}

	var output []chat.Attachment
	for i := range attachments {
		data, err := getAttachmentData(attachments[i].ID)
		if err != nil {
			continue
		}
		output = append(output, chat.Attachment{
			MimeType: attachments[i].Type,
			Data:     data,
		})
	}

	// 4) (optional) warn if we didn't get back as many as we asked for
	if len(output) != len(attachmentIDs) {
		s.log.Warn("mismatch in attachments updated vs requested",
			"requested", len(attachmentIDs),
			"returned", len(output),
		)
	}

	return output, nil
}

// Helper function to encode image to base64
func getAttachmentData(attachmentID string) ([]byte, error) {

	filePath := fmt.Sprintf("data/files/%s", attachmentID)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	return imageData, nil
}
