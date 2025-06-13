package chat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/gorilla/mux"
)

func (s *Service) GetModels(w http.ResponseWriter, r *http.Request) {

	models := s.mr.ListModels()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type ChatListItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Pinned bool   `json:"pinned"`
}

func (s *Service) GetChats(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	chats := make([]ChatListItem, 0)

	rows, err := s.db.Query("SELECT id, title, pinned FROM chats WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var chat ChatListItem
		if err := rows.Scan(&chat.ID, &chat.Title, &chat.Pinned); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type Chat struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Pinned      bool      `json:"pinned"`
	IsStreaming bool      `json:"is_streaming"`
	CreatedAt   int64     `json:"created_at"`
}
type Message struct {
	ID           string       `json:"id"`
	Role         string       `json:"role"`
	Content      string       `json:"content"`
	Attachments  []Attachment `json:"attachments,omitempty"` // Optional field for
	MessageIndex int64        `json:"message_index"`
	CreatedAt    int64        `json:"created_at"`
}

type Attachment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *Service) GetChat(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	id := mux.Vars(r)["id"]

	query := `
        SELECT 
            c.id, c.user_id, c.title, c.model, c.pinned, c.is_streaming, c.created_at,
            m.id, m.role, m.content, m.created_at, m.message_index,
            a.id, a.name, a.type
        FROM chats c
        LEFT JOIN messages m ON c.id = m.chat_id
        LEFT JOIN attachments a ON m.id = a.message_id
        WHERE c.id = ? AND c.user_id = ?
        ORDER BY m.message_index ASC, a.id ASC
    `

	rows, err := s.db.Query(query, id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var chat *Chat
	messages := make(map[string]*Message)

	for rows.Next() {
		var (
			// Chat fields
			cID, cUserID, cTitle, cModel string
			cPinned, cIsStreaming        int
			cCreatedAt                   int64

			// Message fields (nullable)
			mID, mRole, mContent sql.NullString
			mCreatedAt           sql.NullInt64
			mMessageIndex        sql.NullInt64

			// Attachment fields (nullable)
			aID, aName, aType sql.NullString
		)

		err := rows.Scan(
			&cID, &cUserID, &cTitle, &cModel, &cPinned, &cIsStreaming, &cCreatedAt,
			&mID, &mRole, &mContent, &mCreatedAt, &mMessageIndex,
			&aID, &aName, &aType,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Initialize chat on first row
		if chat == nil {
			chat = &Chat{
				ID:          cID,
				UserID:      cUserID,
				Title:       cTitle,
				Model:       cModel,
				Pinned:      cPinned == 1,
				IsStreaming: cIsStreaming == 1,
				CreatedAt:   cCreatedAt,
				Messages:    []Message{},
			}
		}

		// Process message if it exists
		if mID.Valid {
			message, exists := messages[mID.String]
			if !exists {
				message = &Message{
					ID:           mID.String,
					Role:         mRole.String,
					Content:      mContent.String,
					CreatedAt:    mCreatedAt.Int64,
					MessageIndex: mMessageIndex.Int64,
					Attachments:  []Attachment{},
				}
				messages[mID.String] = message
				chat.Messages = append(chat.Messages, *message)
			}

			// Add attachment if it exists
			if aID.Valid {
				attachment := Attachment{
					ID:   aID.String,
					Name: aName.String,
					Type: aType.String,
				}

				// Find the message in chat.Messages and add attachment
				for i := range chat.Messages {
					if chat.Messages[i].ID == mID.String {
						chat.Messages[i].Attachments = append(chat.Messages[i].Attachments, attachment)
						break
					}
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if chat == nil {
		http.Error(w, "Chat not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

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
