package chat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ChatCompletionRequest struct {
	// Message
	Content     string      `json:"content"`
	Attachments []uuid.UUID `json:"attachments,omitempty"`
	// Options
	Model           string `json:"model"`
	ReasoningEffort int32  `json:"reasoning_effort"`
}

func (s *Service) newChat(userID uuid.UUID, request ChatCompletionRequest) (*Chat, error) {

	now := time.Now()
	newChat := Chat{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         fmt.Sprintf("New Chat %d", time.Now().Unix()), // TODO: Generate title based on the first message using a lightweight model
		Model:         request.Model,
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
		return nil, err
	}

	return &newChat, nil

}

func (s *Service) newUserMessage(chatID, userID uuid.UUID, request ChatCompletionRequest) (*chat.Message, error) {

	now := time.Now()
	message := Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		UserID:    userID,
		Role:      "user",
		Status:    "done",
		Model:     request.Model,
		Content:   request.Content,
		CreatedAt: now.UnixMilli(),
		UpdatedAt: now.UnixMilli(),
	}

	_, err := s.db.Exec("INSERT INTO messages (id, chat_id, user_id, stream_id, role, status, model, content, reasoning, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		message.ID, message.ChatID, message.UserID, message.StreamID,
		message.Role, message.Status, message.Model,
		message.Content, message.Reasoning,
		message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		s.log.Warn("failed to insert message into database", "error", err)
		return nil, err
	}

	attachments, err := s.loadAndUpdateAttachments(request.Attachments, message.ID, userID)
	if err != nil {
		s.log.Warn("failed to attach all attachments to message", "message_id", message.ID, "error", err)
		// Note: Consider whether this should be a fatal error or just logged
	}

	return &chat.Message{
		Role:        "user",
		Content:     request.Content,
		Attachments: attachments,
	}, nil

}

func (s *Service) createAssistantMessage(chatID, userID, streamID uuid.UUID, request ChatCompletionRequest, isPremium bool) (uuid.UUID, error) {

	now := time.Now()
	message := Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		UserID:    userID,
		StreamID:  streamID,
		Role:      "assistant",
		Status:    "streaming",
		Model:     request.Model,
		CreatedAt: now.UnixMilli(),
		UpdatedAt: now.UnixMilli(),
	}

	_, err := s.db.Exec("INSERT INTO messages (id, chat_id, user_id, stream_id, role, status, model, content, reasoning, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		message.ID, message.ChatID, message.UserID, message.StreamID,
		message.Role, message.Status, message.Model,
		message.Content, message.Reasoning,
		message.CreatedAt, message.UpdatedAt,
	)
	if err != nil {
		s.log.Warn("failed to insert message into database", "error", err)
		fmt.Println(err)
		return uuid.UUID{}, err
	}

	if isPremium {
		_, err = s.db.Exec("UPDATE user_profile SET usage_premium = usage_premium + 1 WHERE user_id = ?", userID)
	} else {
		_, err = s.db.Exec("UPDATE user_profile SET usage_standard = usage_standard + 1 WHERE user_id = ?", userID)
	}
	if err != nil {
		s.log.Warn("unable to update message usage", "user_id", userID, "error", err)
		fmt.Println("unable to update message usage:", err)
	}

	return message.ID, nil

}

func (s *Service) getChat(chatID, userID uuid.UUID) (*Chat, error) {

	query := `
        SELECT
            c.id, c.user_id, c.title, c.model, c.is_pinned, c.status, c.last_message_at, c.created_at, c.updated_at,
            m.id, m.stream_id, m.role, m.model, m.content, m.reasoning, m.status, m.created_at, m.updated_at,
            a.id, a.name, a.type, a.src, a.created_at
        FROM chats c
        LEFT JOIN messages m ON c.id = m.chat_id
        LEFT JOIN attachments a ON m.id = a.message_id
        WHERE c.id = ? AND c.user_id = ?
        ORDER BY m.created_at ASC, a.created_at ASC
    `

	rows, err := s.db.Query(query, chatID, userID)
	if err != nil {
		s.log.Debug("failed to query chat", "chat_id", chatID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var chat *Chat
	messages := make(map[string]*Message)

	for rows.Next() {
		var (
			// Chat fields
			cID, cUserID, cTitle, cModel, cStatus  string
			cIsPinned                              int
			cLastMessageAt, cCreatedAt, cUpdatedAt int64
			// Message fields (nullable)
			mID, mStreamID, mRole, mModel, mContent, mReasoning, mStatus sql.NullString
			mCreatedAt, mUpdatedAt, aCreatedAt                           sql.NullInt64
			// Attachment fields (nullable)
			aID, aName, aType, aSrc sql.NullString
		)

		err := rows.Scan(
			&cID, &cUserID, &cTitle, &cModel, &cIsPinned, &cStatus, &cLastMessageAt, &cCreatedAt, &cUpdatedAt,
			&mID, &mStreamID, &mRole, &mModel, &mContent, &mReasoning, &mStatus, &mCreatedAt, &mUpdatedAt,
			&aID, &aName, &aType, &aSrc, &aCreatedAt,
		)
		if err != nil {
			s.log.Debug("failed to scan row", "chat_id", chatID, "error", err)
			return nil, err
		}

		// Initialize chat on first row
		if chat == nil {
			chat = &Chat{
				ID:            uuid.MustParse(cID),
				UserID:        uuid.MustParse(cUserID),
				Title:         cTitle,
				Model:         cModel,
				IsPinned:      cIsPinned == 1,
				Status:        cStatus,
				LastMessageAt: cLastMessageAt,
				CreatedAt:     cCreatedAt,
				UpdatedAt:     cUpdatedAt,
				Messages:      []Message{},
			}
		}

		// Process message if it exists
		if mID.Valid {
			message, exists := messages[mID.String]
			if !exists {
				message = &Message{
					ID:          uuid.MustParse(mID.String),
					StreamID:    uuid.MustParse(mStreamID.String),
					Role:        mRole.String,
					Model:       mModel.String,
					Content:     mContent.String,
					Reasoning:   mReasoning.String,
					Status:      mStatus.String,
					CreatedAt:   mCreatedAt.Int64,
					UpdatedAt:   mUpdatedAt.Int64,
					Attachments: []Attachment{},
				}
				messages[mID.String] = message
				chat.Messages = append(chat.Messages, *message)
			}

			// Add attachment if it exists
			if aID.Valid {
				attachment := Attachment{
					ID:        uuid.MustParse(aID.String),
					Name:      aName.String,
					Type:      aType.String,
					Src:       aSrc.String,
					CreatedAt: aCreatedAt.Int64,
				}

				// Find the message in chat.Messages and add attachment
				for i := range chat.Messages {
					if chat.Messages[i].ID == uuid.MustParse(mID.String) { // TODO: that could be a problem
						chat.Messages[i].Attachments = append(chat.Messages[i].Attachments, attachment)
						break
					}
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		s.log.Debug("failed to get rows", "chat_id", chatID, "error", err)
		return nil, err
	}

	if chat == nil {
		return nil, fmt.Errorf("chat not found")
	}

	return chat, nil

}

func (c *Chat) ModelMessages(mf llm.ModelFeatures) ([]*chat.Message, error) {

	var messages []*chat.Message

	for _, msg := range c.Messages {

		message := &chat.Message{
			Role:        msg.Role,
			Content:     msg.Content,
			Attachments: []*chat.Attachment{},
		}

		if mf.HasReasoning {
			message.Reasoning = msg.Reasoning
		}

		for _, att := range msg.Attachments {

			attachment, err := att.ModelAttachment(c.UserID)
			if err != nil {
				fmt.Println("ERROR: ", err)
				continue
			}

			message.Attachments = append(message.Attachments, attachment)

		}

		messages = append(messages, message)

	}

	return messages, nil

}

func (a *Attachment) ModelAttachment(userID uuid.UUID) (*chat.Attachment, error) {

	filePath := fmt.Sprintf("data/users/%s/attachments/%s", userID, a.ID)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	attachmentData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	return &chat.Attachment{
		MimeType: a.Type,
		Data:     attachmentData,
	}, nil

}

func (s *Service) AddMessage(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	// Get chatID from URL params
	chatID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.log.Debug("invalid uuid", "error", err)
		http.Error(w, "invalid_id", http.StatusBadRequest)
		return
	}

	// Decode request body
	var body ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.log.Debug("failed to decode request body", "error", err)
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// Get used model from model router
	model, ok := s.mr.GetModel(body.Model)
	if !ok {
		s.log.Debug("model not supported", "model", body.Model)
		http.Error(w, "model_not_supported", http.StatusBadRequest)
		return
	}

	profile, err := s.getUserProfile(userID)
	if err != nil {
		s.log.Warn("failed to get user profile", "error", err)
		fmt.Println("failed to get user profile:", err)
	}

	if model.Flags.IsPremium {
		if profile.UsagePremium >= profile.LimitPremium {
			s.log.Debug("premium message limit reached", "user_id", userID)
			http.Error(w, "premium_message_limit_reached", http.StatusForbidden)
			return
		}
	} else {
		if profile.UsageStandard >= profile.LimitStandard {
			s.log.Debug("standard message limit reached", "user_id", userID)
			http.Error(w, "standard_message_limit_reached", http.StatusForbidden)
			return
		}
	}

	c, err := s.getChat(chatID, userID)
	if err != nil {
		s.log.Debug("failed to get chat", "chat_id", chatID, "error", err)
		http.Error(w, "get_chat_failed", http.StatusInternalServerError)
		return
	}

	messages, err := c.ModelMessages(model.Features)
	if err != nil {
		s.log.Debug("failed to get chat messages", "chat_id", chatID, "error", err)
		http.Error(w, "get_messages_failed", http.StatusInternalServerError)
		return
	}

	message, err := s.newUserMessage(chatID, userID, body)
	if err != nil {
		s.log.Warn("failed to create user message", "error", err)
		http.Error(w, "create_user_message_failed", http.StatusInternalServerError)
		return
	}

	req := chat.Request{
		Model:               body.Model,
		Temperature:         0.7,
		MaxCompletionTokens: 8192,
		TopP:                1.0,
		Stream:              true,
		ReasoningEffort:     body.ReasoningEffort,
		Stop:                nil,
		Messages:            append(messages, message), // TODO: Consider to split history and new message for better compatibility
		System:              profile.SystemPrompt(),
	}

	compl, err := s.mr.StreamCompletion(req, profile.Options())
	if err != nil {
		s.log.Warn("failed to start a stream", "error", err)
		http.Error(w, "failed to start a stream", http.StatusInternalServerError)
		return
	}

	streamID := uuid.New()
	messageID, err := s.createAssistantMessage(chatID, userID, streamID, body, model.Flags.IsPremium)

	// Add stream to stream pool and return id
	s.sp.Add(streamID.String(), compl)
	compl.OnClose(func(chunk stream.Chunk, serr error) {

		var err error
		if serr != nil {
			s.log.Error("stream failed", "stream_id", streamID, "error", serr)
			_, err = s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
				chunk.Content, chunk.Reasoning, "error", time.Now().UnixMilli(), messageID,
			)
		} else {
			_, err = s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
				chunk.Content, chunk.Reasoning, "done", time.Now().UnixMilli(), messageID,
			)
		}
		if err != nil {
			s.log.Error("storing stream content failed", "stream_id", streamID, "error", err)
			return
		}

	})

	s.log.Debug("stream was started sucessfully", "chat_id", chatID, "stream_id", streamID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"stream_id": streamID,
	}); err != nil {
		s.log.Error("failed to encode response", "error", err)
	}

}

func (s *Service) SendMessage(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
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

	// Get used model from model router
	model, ok := s.mr.GetModel(body.Model)
	if !ok {
		s.log.Debug("model not supported", "model", body.Model)
		http.Error(w, "model_not_supported", http.StatusBadRequest)
		return
	}

	profile, err := s.getUserProfile(userID)
	if err != nil {
		s.log.Warn("failed to get user profile", "error", err)
	}

	if (model.Provider == "anthropic" && profile.AnthropicAPIKey == "") ||
		(model.Provider == "gemini" && profile.GeminiAPIKey == "") ||
		(model.Provider == "ollama" && profile.OllamaBaseURL == "") {
		if model.Flags.IsKeyRequired {
			s.log.Debug("model requires bring your own key", "user_id", userID)
			http.Error(w, "model_requires_key", http.StatusForbidden)
			return
		} else if model.Flags.IsPremium {
			if profile.UsagePremium >= profile.LimitPremium {
				s.log.Debug("premium message limit reached", "user_id", userID)
				http.Error(w, "premium_message_limit_reached", http.StatusForbidden)
				return
			}
		} else {
			if profile.UsageStandard >= profile.LimitStandard {
				s.log.Debug("standard message limit reached", "user_id", userID)
				http.Error(w, "standard_message_limit_reached", http.StatusForbidden)
				return
			}
		}
	}

	c, err := s.newChat(userID, body)
	if err != nil {
		s.log.Warn("failed to create a new chat", "user_id", userID)
		http.Error(w, "create_chat_failed", http.StatusInternalServerError)
		return
	}

	message, err := s.newUserMessage(c.ID, userID, body)
	if err != nil {
		s.log.Warn("failed to create user message", "error", err)
		http.Error(w, "create_user_message_failed", http.StatusInternalServerError)
		return
	}

	req := chat.Request{
		Model:               body.Model,
		Temperature:         0.7,
		MaxCompletionTokens: 8192,
		TopP:                1.0,
		Stream:              true,
		ReasoningEffort:     body.ReasoningEffort,
		Stop:                nil,
		Messages:            []*chat.Message{message}, // TODO: Consider to split history and new message for better compatibility
		System:              profile.SystemPrompt(),
	}

	compl, err := s.mr.StreamCompletion(req, profile.Options())
	if err != nil {
		s.log.Warn("failed to start a stream", "error", err)
		http.Error(w, "failed to start a stream", http.StatusInternalServerError)
		return
	}

	streamID := uuid.New()
	messageID, err := s.createAssistantMessage(c.ID, userID, streamID, body, model.Flags.IsPremium)

	// Add stream to stream pool and return id
	s.sp.Add(streamID.String(), compl)
	compl.OnClose(func(chunk stream.Chunk, serr error) {

		var err error
		if serr != nil {
			s.log.Error("stream failed", "stream_id", streamID, "error", serr)
			_, err = s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
				chunk.Content, chunk.Reasoning, "error", time.Now().UnixMilli(), messageID,
			)
		} else {
			_, err = s.db.Exec("UPDATE messages SET content = ?, reasoning = ?, status = ?, updated_at = ? WHERE id = ?",
				chunk.Content, chunk.Reasoning, "done", time.Now().UnixMilli(), messageID,
			)
		}
		if err != nil {
			s.log.Error("storing stream content failed", "stream_id", streamID, "error", err)
			return
		}

	})

	s.log.Debug("stream was started sucessfully", "chat_id", c.ID, "stream_id", streamID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"chat_id":   c.ID,
		"stream_id": streamID,
	}); err != nil {
		s.log.Error("failed to encode response", "error", err)
	}

}

// LoadAndUpdateAttachments updates message_id on the given attachments and returns their mime_type+data.
func (s *Service) loadAndUpdateAttachments(attachmentIDs []uuid.UUID, messageID, userID uuid.UUID) ([]*chat.Attachment, error) {

	if len(attachmentIDs) == 0 {
		return []*chat.Attachment{}, nil
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
	query := fmt.Sprintf("UPDATE attachments SET message_id = ? WHERE id IN (%s) AND user_id = ? RETURNING id, type", strings.Join(placeholders, ","))

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

	var output []*chat.Attachment
	for i := range attachments {
		data, err := getAttachmentData(userID, attachments[i].ID)
		if err != nil {
			s.log.Error("failed to get attachment data", "error", err)
			continue
		}
		output = append(output, &chat.Attachment{
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
func getAttachmentData(userID, attachmentID uuid.UUID) ([]byte, error) {

	filePath := fmt.Sprintf("data/users/%s/attachments/%s", userID, attachmentID)

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
