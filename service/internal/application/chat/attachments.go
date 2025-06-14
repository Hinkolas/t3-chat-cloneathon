package chat

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Attachment struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id,omitempty"`
	MessageID string `json:"message_id,omitempty"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"`
}

func (s *Service) ListAttachments(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	attachments := make([]Attachment, 0)

	rows, err := s.db.Query("SELECT id, message_id, name, type, created_at FROM attachments WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var attachment Attachment
		if err := rows.Scan(&attachment.ID, &attachment.MessageID, &attachment.Name, &attachment.Type, &attachment.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		attachments = append(attachments, attachment)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attachments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// func (s *Service) GetChat(w http.ResponseWriter, r *http.Request) {

// 	userID := "user-123" // TODO: Replace with context from auth middleware

// 	id := mux.Vars(r)["id"]

// 	query := `
//         SELECT
//             c.id, c.user_id, c.title, c.model, c.is_pinned, c.is_streaming, c.last_message_at, c.created_at, c.updated_at,
//             m.id, m.role, m.model, m.content, m.reasoning, m.created_at, m.updated_at,
//             a.id, a.name, a.type
//         FROM chats c
//         LEFT JOIN messages m ON c.id = m.chat_id
//         LEFT JOIN attachments a ON m.id = a.message_id
//         WHERE c.id = ? AND c.user_id = ?
//         ORDER BY m.created_at ASC, a.id ASC
//     `

// 	rows, err := s.db.Query(query, id, userID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var chat *Chat
// 	messages := make(map[string]*Message)

// 	for rows.Next() {
// 		var (
// 			// Chat fields
// 			cID, cUserID, cTitle, cModel           string
// 			cIsPinned, cIsStreaming                int
// 			cLastMessageAt, cCreatedAt, cUpdatedAt int64
// 			// Message fields (nullable)
// 			mID, mRole, mModel, mContent, mReasoning sql.NullString
// 			mCreatedAt, mUpdatedAt                   sql.NullInt64
// 			// Attachment fields (nullable)
// 			aID, aName, aType sql.NullString
// 		)

// 		err := rows.Scan(
// 			&cID, &cUserID, &cTitle, &cModel, &cIsPinned, &cIsStreaming, &cLastMessageAt, &cCreatedAt, &cUpdatedAt,
// 			&mID, &mRole, &mModel, &mContent, &mReasoning, &mCreatedAt, &mUpdatedAt,
// 			&aID, &aName, &aType,
// 		)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		// Initialize chat on first row
// 		if chat == nil {
// 			chat = &Chat{
// 				ID:            cID,
// 				UserID:        cUserID,
// 				Title:         cTitle,
// 				Model:         cModel,
// 				IsPinned:      cIsPinned == 1,
// 				IsStreaming:   cIsStreaming == 1,
// 				LastMessageAt: cLastMessageAt,
// 				CreatedAt:     cCreatedAt,
// 				UpdatedAt:     cUpdatedAt,
// 				Messages:      []Message{},
// 			}
// 		}

// 		// Process message if it exists
// 		if mID.Valid {
// 			message, exists := messages[mID.String]
// 			if !exists {
// 				message = &Message{
// 					ID:          mID.String,
// 					Role:        mRole.String,
// 					Model:       mModel.String,
// 					Content:     mContent.String,
// 					Reasoning:   mReasoning.String,
// 					CreatedAt:   mCreatedAt.Int64,
// 					UpdatedAt:   mUpdatedAt.Int64,
// 					Attachments: []Attachment{},
// 				}
// 				messages[mID.String] = message
// 				chat.Messages = append(chat.Messages, *message)
// 			}

// 			// Add attachment if it exists
// 			if aID.Valid {
// 				attachment := Attachment{
// 					ID:   aID.String,
// 					Name: aName.String,
// 					Type: aType.String,
// 				}

// 				// Find the message in chat.Messages and add attachment
// 				for i := range chat.Messages {
// 					if chat.Messages[i].ID == mID.String {
// 						chat.Messages[i].Attachments = append(chat.Messages[i].Attachments, attachment)
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}

// 	if err := rows.Err(); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if chat == nil {
// 		http.Error(w, "Chat not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(chat); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// }

func (s *Service) DeleteAttachment(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	id := mux.Vars(r)["id"]

	// Delete the chat (this will cascade delete messages and attachments if foreign keys are set up properly)
	result, err := s.db.Exec("DELETE FROM attachments WHERE id = ? AND user_id = ?", id, userID)
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
		http.Error(w, "Attachment not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// type PatchChatRequest struct {
// 	Title    *string `json:"title,omitempty"`
// 	IsPinned *bool   `json:"is_pinned,omitempty"`
// }

// func (s *Service) EditChat(w http.ResponseWriter, r *http.Request) {

// 	userID := "user-123" // TODO: Replace with context from auth middleware

// 	id := mux.Vars(r)["id"]

// 	var req PatchChatRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Title == nil && req.IsPinned == nil {
// 		http.Error(w, "No fields to update", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Title != nil {
// 		result, err := s.db.Exec("UPDATE chats SET title = ? WHERE id = ? AND user_id = ?", *req.Title, id, userID)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		rowsAffected, err := result.RowsAffected()
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		if rowsAffected == 0 {
// 			http.Error(w, "Chat not found", http.StatusNotFound)
// 			return
// 		}
// 	}

// 	if req.IsPinned != nil {
// 		result, err := s.db.Exec("UPDATE chats SET is_pinned = ? WHERE id = ? AND user_id = ?", *req.IsPinned, id, userID)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		rowsAffected, err := result.RowsAffected()
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		if rowsAffected == 0 {
// 			http.Error(w, "Chat not found", http.StatusNotFound)
// 			return
// 		}
// 	}

// 	w.WriteHeader(http.StatusNoContent)

// }
