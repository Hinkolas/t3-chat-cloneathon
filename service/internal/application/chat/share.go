package chat

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type SharedChat struct {
	ID       string          `json:"id"`
	UserID   string          `json:"user_id"`
	Title    string          `json:"title"`
	Model    string          `json:"model"`
	Messages []SharedMessage `json:"messages"`
}

type SharedMessage struct {
	ID          string             `json:"id"`
	Role        string             `json:"role"`
	Model       string             `json:"model"`
	Content     string             `json:"content"`
	Reasoning   string             `json:"reasoning"`
	Attachments []SharedAttachment `json:"attachments"`
}

type SharedAttachment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Src  string `json:"src"`
}

func (s *Service) GetSharedChat(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	query := `
        SELECT
            c.id, c.user_id, c.title, c.model,
            m.id, m.role, m.model, m.content, m.reasoning,
            a.id, a.name, a.type, a.src
        FROM chats c
        LEFT JOIN messages m ON c.id = m.chat_id
        LEFT JOIN attachments a ON m.id = a.message_id
        WHERE c.id = ? AND c.created_at <= c.shared_at AND m.created_at <= c.shared_at
        ORDER BY m.created_at ASC, a.created_at ASC
    `

	rows, err := s.db.Query(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var chat *SharedChat
	messages := make(map[string]*SharedMessage)

	for rows.Next() {
		var (
			// Chat fields
			cID, cUserID, cTitle, cModel string
			// Message fields (nullable)
			mID, mRole, mModel, mContent, mReasoning sql.NullString
			// Attachment fields (nullable)
			aID, aName, aType, aSrc sql.NullString
		)

		err := rows.Scan(
			&cID, &cUserID, &cTitle, &cModel,
			&mID, &mRole, &mModel, &mContent, &mReasoning,
			&aID, &aName, &aType, &aSrc,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Initialize chat on first row
		if chat == nil {
			chat = &SharedChat{
				ID:       cID,
				UserID:   cUserID,
				Title:    cTitle,
				Model:    cModel,
				Messages: []SharedMessage{},
			}
		}

		// Process message if it exists
		if mID.Valid {
			message, exists := messages[mID.String]
			if !exists {
				message = &SharedMessage{
					ID:          mID.String,
					Role:        mRole.String,
					Model:       mModel.String,
					Content:     mContent.String,
					Reasoning:   mReasoning.String,
					Attachments: []SharedAttachment{},
				}
				messages[mID.String] = message
				chat.Messages = append(chat.Messages, *message)
			}

			// Add attachment if it exists
			if aID.Valid {
				attachment := SharedAttachment{
					ID:   aID.String,
					Name: aName.String,
					Type: aType.String,
					Src:  aSrc.String,
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
