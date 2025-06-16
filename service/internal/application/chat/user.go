package chat

import (
	"encoding/json"
	"net/http"
	"strings"
)

type UserProfile struct {
	UserID string `db:"user_id" json:"user_id"`
	// Provider Options
	AnthropicAPIKey string `db:"anthropic_api_key" json:"anthropic_api_key"`
	OpenAIAPIKey    string `db:"openai_api_key" json:"openai_api_key"`
	GeminiAPIKey    string `db:"gemini_api_key" json:"gemini_api_key"`
	OllamaBaseURL   string `db:"ollama_base_url" json:"ollama_base_url"`
	// Customization
	CustomUserName       string `db:"custom_user_name" json:"custom_user_name"`
	CustomUserProfession string `db:"custom_user_profession" json:"custom_user_profession"`
	CustomAssistantTrait string `db:"custom_assistant_trait" json:"custom_assistant_trait"`
	CustomContext        string `db:"custom_context" json:"custom_context"`
}

func (s *Service) getUserProfile(userID string) (*UserProfile, error) {
	profile := &UserProfile{}
	err := s.db.QueryRow("SELECT user_id, anthropic_api_key, openai_api_key, gemini_api_key, ollama_base_url, custom_user_name, custom_user_profession, custom_assistant_trait, custom_context FROM user_profile WHERE user_id = ?", userID).Scan(
		&profile.UserID, &profile.AnthropicAPIKey, &profile.OpenAIAPIKey, &profile.GeminiAPIKey, &profile.OllamaBaseURL, &profile.CustomUserName, &profile.CustomUserProfession, &profile.CustomAssistantTrait, &profile.CustomContext)
	return profile, err
}

func (s *Service) PatchUserProfile(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	var updates UserProfile
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var setParts []string
	var args []any

	if updates.AnthropicAPIKey != "" {
		setParts = append(setParts, "anthropic_api_key = ?")
		args = append(args, updates.AnthropicAPIKey)
	}
	if updates.OpenAIAPIKey != "" {
		setParts = append(setParts, "openai_api_key = ?")
		args = append(args, updates.OpenAIAPIKey)
	}
	if updates.GeminiAPIKey != "" {
		setParts = append(setParts, "gemini_api_key = ?")
		args = append(args, updates.GeminiAPIKey)
	}
	if updates.OllamaBaseURL != "" {
		setParts = append(setParts, "ollama_base_url = ?")
		args = append(args, updates.OllamaBaseURL)
	}
	if updates.CustomUserName != "" {
		setParts = append(setParts, "custom_user_name = ?")
		args = append(args, updates.CustomUserName)
	}
	if updates.CustomUserProfession != "" {
		setParts = append(setParts, "custom_user_profession = ?")
		args = append(args, updates.CustomUserProfession)
	}
	if updates.CustomAssistantTrait != "" {
		setParts = append(setParts, "custom_assistant_trait = ?")
		args = append(args, updates.CustomAssistantTrait)
	}
	if updates.CustomContext != "" {
		setParts = append(setParts, "custom_context = ?")
		args = append(args, updates.CustomContext)
	}

	if len(setParts) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	query := "UPDATE user_profile SET " + strings.Join(setParts, ", ") + " WHERE user_id = ?"
	args = append(args, userID)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (s *Service) GetUserProfile(w http.ResponseWriter, r *http.Request) {

	userID := "user-123" // TODO: Replace with context from auth middleware

	profile, err := s.getUserProfile(userID)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
