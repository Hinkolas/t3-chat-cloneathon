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

	// TODO: Replace with a more elegant solution
	// Update attachments with the message ID
	if len(body.Attachments) > 0 {
		err = s.updateAttachmentsMessageID(body.Attachments, message.ID, userID)
		if err != nil {
			s.log.Warn("failed to update attachments with message ID", "error", err)
			// Note: Consider whether this should be a fatal error or just logged
		}
	}

	profile, err := s.getUserProfile(userID)
	if err != nil {
		s.log.Warn("failed to get user profile", "error", err)
	}

	prompt := fmt.Sprintf(`You are a personalized AI assistant. Your primary goal is to be helpful, engaging, and tailored to the user's specific needs and preferences.

# Your Persona
Your personality should be guided by the following traits: %s.
Embody these traits in your language, tone, and the structure of your responses. If no traits are specified, adopt a generally helpful, friendly, and curious persona.

# User Information
You have been provided with the following information about the user. Use it to personalize the conversation naturally.

- **User's Name:** %s
  - Address the user by this name occasionally to create a more personal connection. Use it in greetings or when it feels natural, but avoid overusing it. If no name is provided, use neutral and friendly greetings.

- **User's Role:** %s
  - The user identifies as a %s. Keep this role in mind to understand their perspective. Your examples, analogies, and suggestions may be more effective if they relate to this role.

- **User's Context & Preferences:** %s
  - This is the most important information for personalization. It contains the user's interests, values, and preferences.
  - Proactively use this context to tailor your responses. For example, if the user is interested in "philosophy," you can frame answers or provide examples through a philosophical lens. If they prefer "concise" answers, keep your responses direct and to the point.

# Core Instructions
1.  **Integrate, Don't Recite:** Do not state the user's information back to them (e.g., "As an engineer, you might like..."). Instead, let the information subtly guide your word choice, examples, and tone.
2.  **Prioritize User Context:** The user's context in '%s' should be the primary driver of your personalization.
3.  **Maintain Your Persona:** Consistently apply the traits from '%s' throughout the entire conversation.
4.  **Be Helpful:** Your ultimate purpose is to assist the user effectively with their requests.`,
		profile.CustomAssistantTrait,
		profile.CustomUserName,
		profile.CustomUserProfession,
		profile.CustomUserProfession,
		profile.CustomContext,
		profile.CustomContext,        // Add this line
		profile.CustomAssistantTrait, // Add this line
	)

	compl, err := s.mr.StreamCompletion(req, chat.Options{
		AnthropicAPIKey: profile.AnthropicAPIKey,
		GeminiAPIKey:    profile.GeminiAPIKey,
		OpenAIAPIKey:    profile.OpenAIAPIKey,
		OllamaBaseURL:   profile.OllamaBaseURL,
		SystemPrompt:    prompt,
	})
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

	messages := []chat.Message{
		{
			Role:    "user",
			Content: body.Content,
		},
	}

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

	// TODO: Replace with a more elegant solution
	// Update attachments with the message ID
	if len(body.Attachments) > 0 {
		err = s.updateAttachmentsMessageID(body.Attachments, message.ID, userID)
		if err != nil {
			s.log.Warn("failed to update attachments with message ID", "error", err)
			// Note: Consider whether this should be a fatal error or just logged
		}
	}

	profile, err := s.getUserProfile(userID)
	if err != nil {
		s.log.Warn("failed to get user profile", "error", err)
	}

	prompt := fmt.Sprintf(`You are a personalized AI assistant. Your primary goal is to be helpful, engaging, and tailored to the user's specific needs and preferences.

# Your Persona
Your personality should be guided by the following traits: %s.
Embody these traits in your language, tone, and the structure of your responses. If no traits are specified, adopt a generally helpful, friendly, and curious persona.

# User Information
You have been provided with the following information about the user. Use it to personalize the conversation naturally.

- **User's Name:** %s
  - Address the user by this name occasionally to create a more personal connection. Use it in greetings or when it feels natural, but avoid overusing it. If no name is provided, use neutral and friendly greetings.

- **User's Role:** %s
  - The user identifies as a %s. Keep this role in mind to understand their perspective. Your examples, analogies, and suggestions may be more effective if they relate to this role.

- **User's Context & Preferences:** %s
  - This is the most important information for personalization. It contains the user's interests, values, and preferences.
  - Proactively use this context to tailor your responses. For example, if the user is interested in "philosophy," you can frame answers or provide examples through a philosophical lens. If they prefer "concise" answers, keep your responses direct and to the point.

# Core Instructions
1.  **Integrate, Don't Recite:** Do not state the user's information back to them (e.g., "As an engineer, you might like..."). Instead, let the information subtly guide your word choice, examples, and tone.
2.  **Prioritize User Context:** The user's context in '%s' should be the primary driver of your personalization.
3.  **Maintain Your Persona:** Consistently apply the traits from '%s' throughout the entire conversation.
4.  **Be Helpful:** Your ultimate purpose is to assist the user effectively with their requests.`,
		profile.CustomAssistantTrait,
		profile.CustomUserName,
		profile.CustomUserProfession,
		profile.CustomUserProfession,
		profile.CustomContext,
		profile.CustomContext,        // Add this line
		profile.CustomAssistantTrait, // Add this line
	)

	compl, err := s.mr.StreamCompletion(req, chat.Options{
		AnthropicAPIKey: profile.AnthropicAPIKey,
		GeminiAPIKey:    profile.GeminiAPIKey,
		OpenAIAPIKey:    profile.OpenAIAPIKey,
		OllamaBaseURL:   profile.OllamaBaseURL,
		SystemPrompt:    prompt,
	})
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

// Helper method to update attachments
func (s *Service) updateAttachmentsMessageID(attachmentIDs []string, messageID, userID string) error {
	if len(attachmentIDs) == 0 {
		return nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(attachmentIDs))
	args := make([]interface{}, 0, len(attachmentIDs)+2)

	args = append(args, messageID)
	for i := range attachmentIDs {
		placeholders[i] = "?"
		args = append(args, attachmentIDs[i])
	}
	args = append(args, userID)

	query := fmt.Sprintf("UPDATE attachments SET message_id = ? WHERE id IN (%s) AND user_id = ? AND message_id = ''",
		strings.Join(placeholders, ","))

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update attachments: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if int(rowsAffected) != len(attachmentIDs) {
		s.log.Warn("not all attachments were updated",
			"expected", len(attachmentIDs),
			"updated", rowsAffected)
	}

	return nil
}
