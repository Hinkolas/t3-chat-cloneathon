package chat

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Service) ListModels(w http.ResponseWriter, r *http.Request) {

	models := s.mr.ListModels()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type ChatListItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	IsPinned      bool   `json:"is_pinned"`
	LastMessageAt int64  `json:"last_message_at"`
	CreatedAt     int64  `json:"created_at"`
}

func (s *Service) ListChats(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	chats := make([]ChatListItem, 0)

	rows, err := s.db.Query("SELECT id, title, is_pinned, last_message_at, created_at FROM chats WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var chat ChatListItem
		if err := rows.Scan(&chat.ID, &chat.Title, &chat.IsPinned, &chat.LastMessageAt, &chat.CreatedAt); err != nil {
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
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Model       string `json:"model"`
	IsPinned    bool   `json:"is_pinned"`
	IsStreaming bool   `json:"is_streaming"`

	LastMessageAt int64 `json:"last_message_at"`
	CreatedAt     int64 `json:"created_at"`
	UpdatedAt     int64 `json:"updated_at"`

	Messages []Message `json:"messages"`
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
            c.id, c.user_id, c.title, c.model, c.is_pinned, c.is_streaming, c.last_message_at, c.created_at, c.updated_at,
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
			cID, cUserID, cTitle, cModel           string
			cIsPinned, cIsStreaming                int
			cLastMessageAt, cCreatedAt, cUpdatedAt int64
			// Message fields (nullable)
			mID, mRole, mContent sql.NullString
			mCreatedAt           sql.NullInt64
			mMessageIndex        sql.NullInt64

			// Attachment fields (nullable)
			aID, aName, aType sql.NullString
		)

		err := rows.Scan(
			&cID, &cUserID, &cTitle, &cModel, &cIsPinned, &cIsStreaming, &cLastMessageAt, &cCreatedAt, &cUpdatedAt,
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
				ID:            cID,
				UserID:        cUserID,
				Title:         cTitle,
				Model:         cModel,
				IsPinned:      cIsPinned == 1,
				IsStreaming:   cIsStreaming == 1,
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
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *Service) DeleteChat(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	id := mux.Vars(r)["id"]

	// Delete the chat (this will cascade delete messages and attachments if foreign keys are set up properly)
	result, err := s.db.Exec("DELETE FROM chats WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

type PatchChatRequest struct {
	Title    *string `json:"title,omitempty"`
	IsPinned *bool   `json:"is_pinned,omitempty"`
}

func (s *Service) EditChat(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	id := mux.Vars(r)["id"]

	var req PatchChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == nil && req.IsPinned == nil {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	if req.Title != nil {
		result, err := s.db.Exec("UPDATE chats SET title = ? WHERE id = ? AND user_id = ?", *req.Title, id, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
	}

	if req.IsPinned != nil {
		result, err := s.db.Exec("UPDATE chats SET is_pinned = ? WHERE id = ? AND user_id = ?", *req.IsPinned, id, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)

}
