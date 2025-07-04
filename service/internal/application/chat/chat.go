package chat

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Chat struct {
	ID     uuid.UUID `json:"id,omitzero"`
	UserID uuid.UUID `json:"user_id,omitzero"`

	Title         string `json:"title"`
	Model         string `json:"model"`
	IsPinned      bool   `json:"is_pinned"`
	Status        string `json:"status"` // e.g. "streaming", "done", "error"
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	LastMessageAt int64  `json:"last_message_at"`
	SharedAt      int64  `json:"shared_at"`

	Messages []Message `json:"messages"`
}

type Message struct {
	ID       uuid.UUID `json:"id,omitzero"`
	ChatID   uuid.UUID `json:"chat_id,omitzero"`
	UserID   uuid.UUID `json:"user_id,omitzero"`
	StreamID uuid.UUID `json:"stream_id,omitzero"`

	Role      string `json:"role"`
	Model     string `json:"model"`
	Content   string `json:"content"`
	Reasoning string `json:"reasoning,omitempty"`

	Status    string `json:"status"` // e.g. "streaming", "done", "error"
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	Attachments []Attachment `json:"attachments"`
}

type ChatListItem struct {
	ID            uuid.UUID `json:"id,omitzero"`
	Title         string    `json:"title"`
	IsPinned      bool      `json:"is_pinned"`
	Status        string    `json:"status"`
	LastMessageAt int64     `json:"last_message_at"`
	CreatedAt     int64     `json:"created_at"`
	SharedAt      int64     `json:"shared_at"`
}

type PatchChatRequest struct {
	Title    *string `json:"title,omitempty"`
	IsPinned *bool   `json:"is_pinned,omitempty"`
	Model    *string `json:"model,omitempty"`
	SharedAt *int64  `json:"shared_at,omitempty"`
}

func (s *Service) ListChats(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	chats := make([]ChatListItem, 0)

	rows, err := s.db.Query("SELECT id, title, is_pinned, status, last_message_at, created_at, shared_at FROM chats WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var chat ChatListItem
		if err := rows.Scan(&chat.ID, &chat.Title, &chat.IsPinned, &chat.Status, &chat.LastMessageAt, &chat.CreatedAt, &chat.SharedAt); err != nil {
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

func (s *Service) GetChat(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

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

func (s *Service) EditChat(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	var req PatchChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == nil && req.IsPinned == nil && req.SharedAt == nil && req.Model == nil {
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

	if req.Model != nil {
		_, ok := s.mr.GetModel(*req.Model)
		if ok {

			result, err := s.db.Exec("UPDATE chats SET model = ? WHERE id = ? AND user_id = ?", *req.Model, id, userID)
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
	}

	if req.SharedAt != nil {

		result, err := s.db.Exec("UPDATE chats SET shared_at = ? WHERE id = ? AND user_id = ?", *req.SharedAt, id, userID)
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
