package chat

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

//go:embed system.txt
var systemTemplate string

type UserProfile struct {
	UserID uuid.UUID `json:"user_id,omitzero"`
	// User
	Username string `json:"username"` // Username of the user
	Email    string `json:"email"`    // Email of the user
	// Usage
	LimitStandard int32 `json:"limit_standard"`
	LimitPremium  int32 `json:"limit_premium"`
	UsageStandard int32 `json:"usage_standard"`
	UsagePremium  int32 `json:"usage_premium"`
	// Provider Options
	AnthropicAPIKey string `json:"anthropic_api_key"`
	OpenAIAPIKey    string `json:"openai_api_key"`
	GeminiAPIKey    string `json:"gemini_api_key"`
	OllamaBaseURL   string `json:"ollama_base_url"`
	// Customization
	CustomUserName       string `json:"custom_user_name"`
	CustomUserProfession string `json:"custom_user_profession"`
	CustomAssistantTrait string `json:"custom_assistant_trait"`
	CustomContext        string `json:"custom_context"`
}

func (p *UserProfile) SystemPrompt() string {

	tmpl, err := template.New("system").Parse(systemTemplate)
	if err != nil {
		fmt.Println("failed to parse system template:", err)
		return systemTemplate // fallback to original template on parse error
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, p)
	if err != nil {
		fmt.Println("failed to execute system template:", err)
		return systemTemplate // fallback to original template on execution error
	}

	return buf.String()

}

func (p *UserProfile) Options() map[string]string {
	options := make(map[string]string)
	if p.AnthropicAPIKey != "" {
		options["anthropic_api_key"] = p.AnthropicAPIKey
	}
	if p.OpenAIAPIKey != "" {
		options["openai_api_key"] = p.OpenAIAPIKey
	}
	if p.GeminiAPIKey != "" {
		options["gemini_api_key"] = p.GeminiAPIKey
	}
	if p.OllamaBaseURL != "" {
		options["ollama_base_url"] = p.OllamaBaseURL
	}
	return options
}

func (s *Service) getUserProfile(userID uuid.UUID) (*UserProfile, error) {
	profile := &UserProfile{}
	err := s.db.QueryRow(`
		SELECT up.user_id, u.username, u.email, up.limit_standard, up.limit_premium, up.usage_standard, up.usage_premium,
		       up.anthropic_api_key, up.openai_api_key, up.gemini_api_key, up.ollama_base_url,
		       up.custom_user_name, up.custom_user_profession, up.custom_assistant_trait, up.custom_context
		FROM user_profile up
		JOIN users u ON up.user_id = u.id
		WHERE up.user_id = ?`, userID).Scan(
		&profile.UserID, &profile.Username, &profile.Email, &profile.LimitStandard, &profile.LimitPremium, &profile.UsageStandard, &profile.UsagePremium, &profile.AnthropicAPIKey, &profile.OpenAIAPIKey, &profile.GeminiAPIKey, &profile.OllamaBaseURL, &profile.CustomUserName, &profile.CustomUserProfession, &profile.CustomAssistantTrait, &profile.CustomContext)
	return profile, err
}

type PatchProfileRequest struct {
	AnthropicAPIKey      *string `json:"anthropic_api_key,omitempty"`
	OpenAIAPIKey         *string `json:"openai_api_key,omitempty"`
	GeminiAPIKey         *string `json:"gemini_api_key,omitempty"`
	OllamaBaseURL        *string `json:"ollama_base_url,omitempty"`
	CustomUserName       *string `json:"custom_user_name,omitempty"`
	CustomUserProfession *string `json:"custom_user_profession,omitempty"`
	CustomAssistantTrait *string `json:"custom_assistant_trait,omitempty"`
	CustomContext        *string `json:"custom_context,omitempty"`
}

func (s *Service) UpsertUserProfile(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	var updates PatchProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// we always have the PK:
	cols := []string{"user_id"}
	placeholder := []string{"?"}
	values := []any{userID}

	// for the ON CONFLICT ... DO UPDATE clause:
	var updateSets []string

	// helper to append a col if the field is non‐nil
	appendField := func(col string, v *string) {
		if v != nil {
			cols = append(cols, col)
			placeholder = append(placeholder, "?")
			values = append(values, *v)
			// EXCLUDED.col references the incoming value
			updateSets = append(updateSets, col+" = EXCLUDED."+col)
		}
	}

	appendField("anthropic_api_key", updates.AnthropicAPIKey)
	appendField("openai_api_key", updates.OpenAIAPIKey)
	appendField("gemini_api_key", updates.GeminiAPIKey)
	appendField("ollama_base_url", updates.OllamaBaseURL)
	appendField("custom_user_name", updates.CustomUserName)
	appendField("custom_user_profession", updates.CustomUserProfession)
	appendField("custom_assistant_trait", updates.CustomAssistantTrait)
	appendField("custom_context", updates.CustomContext)

	if len(updateSets) == 0 {
		http.Error(w, "No fields to upsert", http.StatusBadRequest)
		return
	}

	// build the final upsert SQL
	// INSERT INTO user_profile (col1, col2, …)
	// VALUES (?, ?, …)
	// ON CONFLICT (user_id) DO UPDATE SET col2 = EXCLUDED.col2, …
	query := fmt.Sprintf(
		"INSERT INTO user_profile (%s) VALUES (%s) ON CONFLICT(user_id) DO UPDATE SET %s",
		strings.Join(cols, ", "),
		strings.Join(placeholder, ", "),
		strings.Join(updateSets, ", "),
	)

	if _, err := s.db.Exec(query, values...); err != nil {
		http.Error(w, "Failed to upsert profile: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})

}

func (s *Service) GetUserProfile(w http.ResponseWriter, r *http.Request) {

	// Get userID from auth middleware, ok if authenticated
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		s.log.Debug("User is not authenticated")
		http.Error(w, "not_authenticated", http.StatusUnauthorized)
		return
	}

	profile, err := s.getUserProfile(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UserProfile{UserID: userID})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)

}
