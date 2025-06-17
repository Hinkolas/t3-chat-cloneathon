package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Attachment struct {
	ID        uuid.UUID `json:"id,omitzero"`
	UserId    uuid.UUID `json:"user_id,omitzero"`
	MessageID uuid.UUID `json:"message_id,omitzero"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Src       string    `json:"src"`
	CreatedAt int64     `json:"created_at"`
}

func (s *Service) ListAttachments(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	attachments := make([]Attachment, 0)

	rows, err := s.db.Query("SELECT id, message_id, name, type, src, created_at FROM attachments WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var attachment Attachment
		if err := rows.Scan(&attachment.ID, &attachment.MessageID, &attachment.Name, &attachment.Type, &attachment.Src, &attachment.CreatedAt); err != nil {
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

func (s *Service) GetAttachment(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	var attachment Attachment
	err := s.db.QueryRow("SELECT id, message_id, name, type, src, created_at FROM attachments WHERE id = ? AND user_id = ?", id, userID).Scan(
		&attachment.ID, &attachment.MessageID, &attachment.Name, &attachment.Type, &attachment.Src, &attachment.CreatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Attachment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", attachment.Type)

	http.ServeFile(w, r, fmt.Sprintf("data/files/%s", attachment.ID))

}

func (s *Service) UploadAttachment(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create attachment record
	now := time.Now()
	attachment := Attachment{
		// ID:        fmt.Sprintf("att_%d", now.UnixNano()),
		ID:        uuid.New(),
		UserId:    userID,
		Name:      header.Filename,
		Type:      header.Header.Get("Content-Type"),
		Src:       fmt.Sprintf("http://127.0.0.1:3141/v1/attachments/att_%d/", now.UnixNano()), // TODO: Replace with proper location
		CreatedAt: now.UnixMilli(),
	}

	// Save to database
	_, err = s.db.Exec("INSERT INTO attachments (id, user_id, message_id, name, type, src, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		attachment.ID, attachment.UserId, attachment.MessageID, attachment.Name, attachment.Type, attachment.Src, attachment.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create attachments directory if it doesn't exist
	os.MkdirAll("data/files", 0755)

	// Save file to disk
	dst, err := os.Create(fmt.Sprintf("data/files/att_%d", now.UnixNano()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(attachment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *Service) DeleteAttachment(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

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
